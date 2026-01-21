package cloud189

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/cloudreve/Cloudreve/v4/ent"
	"github.com/cloudreve/Cloudreve/v4/pkg/boolset"
	"github.com/cloudreve/Cloudreve/v4/pkg/conf"
	"github.com/cloudreve/Cloudreve/v4/pkg/filemanager/driver"
	"github.com/cloudreve/Cloudreve/v4/pkg/filemanager/fs"
	"github.com/cloudreve/Cloudreve/v4/pkg/logging"
	"github.com/cloudreve/Cloudreve/v4/pkg/request"
	"github.com/go-resty/resty/v2"
)

const (
	// DefaultChunkSize 默认分片大小 10MB
	DefaultChunkSize int64 = 10485760
	// RootFolderID 根目录ID
	RootFolderID = "-11"
)

var (
	capabilities = &driver.Capabilities{
		StaticFeatures: &boolset.BooleanSet{},
	}
)

func init() {
	boolset.Sets(map[driver.HandlerCapability]bool{
		driver.HandlerCapabilityProxyRequired: true,
	}, capabilities.StaticFeatures)
}

// Driver 天翼云盘驱动
type Driver struct {
	policy     *ent.StoragePolicy
	client     *resty.Client
	httpClient request.Client
	l          logging.Logger
	config     conf.ConfigProvider

	username   string
	password   string
	sessionKey string
	rsa        Rsa
}

// New 创建天翼云盘驱动实例
func New(policy *ent.StoragePolicy, l logging.Logger, config conf.ConfigProvider) (*Driver, error) {
	// 从策略配置中获取用户名和密码
	username := policy.AccessKey
	password := policy.SecretKey

	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	client := resty.New().
		SetHeader("Referer", "https://cloud.189.cn/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	driver := &Driver{
		policy:     policy,
		client:     client,
		httpClient: request.NewClient(config, request.WithLogger(l)),
		l:          l,
		config:     config,
		username:   username,
		password:   password,
	}

	// 执行登录
	if err := driver.newLogin(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return driver, nil
}

// List 列出目录下的文件
func (d *Driver) List(ctx context.Context, base string, onProgress driver.ListProgressFunc, recursive bool) ([]fs.PhysicalObject, error) {
	res := make([]fs.PhysicalObject, 0)
	
	// 天翼云盘使用文件夹ID而不是路径
	folderID := RootFolderID
	if base != "" && base != "/" {
		// 这里需要根据路径查找文件夹ID，简化处理
		folderID = base
	}

	pageNum := 1
	for {
		var resp Files
		_, err := d.request("https://cloud.189.cn/api/open/file/listFiles.action", http.MethodGet, func(req *resty.Request) {
			req.SetQueryParams(map[string]string{
				"pageSize":   "60",
				"pageNum":    strconv.Itoa(pageNum),
				"mediaType":  "0",
				"folderId":   folderID,
				"iconOption": "5",
				"orderBy":    "lastOpTime",
				"descending": "true",
			})
		}, &resp)
		if err != nil {
			return nil, err
		}

		if resp.FileListAO.Count == 0 {
			break
		}

		// 处理文件夹
		for _, folder := range resp.FileListAO.FolderList {
			lastOpTime, _ := parseCNTime(folder.LastOpTime)
			res = append(res, fs.PhysicalObject{
				Name:         folder.Name,
				RelativePath: folder.Name,
				Source:       strconv.FormatInt(folder.Id, 10),
				Size:         0,
				IsDir:        true,
				LastModify:   lastOpTime,
			})
		}

		// 处理文件
		for _, file := range resp.FileListAO.FileList {
			lastOpTime, _ := parseCNTime(file.LastOpTime)
			res = append(res, fs.PhysicalObject{
				Name:         file.Name,
				RelativePath: file.Name,
				Source:       strconv.FormatInt(file.Id, 10),
				Size:         file.Size,
				IsDir:        false,
				LastModify:   lastOpTime,
			})
		}

		onProgress(len(res))
		pageNum++
	}

	return res, nil
}

// Open 打开文件（不支持）
func (d *Driver) Open(ctx context.Context, path string) (*os.File, error) {
	return nil, errors.New("not implemented")
}

// LocalPath 获取本地路径（不支持）
func (d *Driver) LocalPath(ctx context.Context, path string) string {
	return ""
}

// Put 上传文件
func (d *Driver) Put(ctx context.Context, file *fs.UploadRequest) error {
	defer file.Close()

	// 获取session key
	sessionKey, err := d.getSessionKey()
	if err != nil {
		return err
	}
	d.sessionKey = sessionKey

	// 获取父文件夹ID
	parentFolderID := RootFolderID
	if file.Props.SavePath != "" {
		dir := path.Dir(file.Props.SavePath)
		if dir != "" && dir != "/" && dir != "." {
			parentFolderID = dir
		}
	}

	// 计算分片数量
	chunkSize := DefaultChunkSize
	count := int64(math.Ceil(float64(file.Props.Size) / float64(chunkSize)))

	// 初始化分片上传
	res, err := d.uploadRequest("/person/initMultiUpload", map[string]string{
		"parentFolderId": parentFolderID,
		"fileName":       encode(file.Props.Uri.Name()),
		"fileSize":       strconv.FormatInt(file.Props.Size, 10),
		"sliceSize":      strconv.FormatInt(chunkSize, 10),
		"lazyCheck":      "1",
	}, nil)
	if err != nil {
		return err
	}

	var uploadData map[string]interface{}
	if err := json.Unmarshal(res, &uploadData); err != nil {
		return err
	}

	data, ok := uploadData["data"].(map[string]interface{})
	if !ok {
		return errors.New("invalid upload response")
	}

	uploadFileId, ok := data["uploadFileId"].(string)
	if !ok {
		return errors.New("uploadFileId not found")
	}

	// 上传分片
	var finish int64 = 0
	md5s := make([]string, 0)
	md5Sum := md5.New()

	for i := int64(1); i <= count; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		byteSize := file.Props.Size - finish
		if chunkSize < byteSize {
			byteSize = chunkSize
		}

		byteData := make([]byte, byteSize)
		n, err := io.ReadFull(file, byteData)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}

		finish += int64(n)
		md5Bytes := getMd5(byteData)
		md5Hex := hex.EncodeToString(md5Bytes)
		md5Base64 := base64.StdEncoding.EncodeToString(md5Bytes)
		md5s = append(md5s, strings.ToUpper(md5Hex))
		md5Sum.Write(byteData)

		// 获取上传URL
		var urlResp UploadUrlsResp
		_, err = d.uploadRequest("/person/getMultiUploadUrls", map[string]string{
			"partInfo":     fmt.Sprintf("%s-%s", strconv.FormatInt(i, 10), md5Base64),
			"uploadFileId": uploadFileId,
		}, &urlResp)
		if err != nil {
			return err
		}

		uploadData := urlResp.UploadUrls["partNumber_"+strconv.FormatInt(i, 10)]
		requestURL := uploadData.RequestURL
		uploadHeaders := strings.Split(decodeURIComponent(uploadData.RequestHeader), "&")

		// 上传分片
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL, bytes.NewReader(byteData))
		if err != nil {
			return err
		}

		for _, v := range uploadHeaders {
			idx := strings.Index(v, "=")
			if idx > 0 {
				req.Header.Set(v[0:idx], v[idx+1:])
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
	}

	// 完成上传
	fileMd5 := hex.EncodeToString(md5Sum.Sum(nil))
	sliceMd5 := fileMd5
	if file.Props.Size > chunkSize {
		h := md5.New()
		h.Write([]byte(strings.Join(md5s, "\n")))
		sliceMd5 = hex.EncodeToString(h.Sum(nil))
	}

	_, err = d.uploadRequest("/person/commitMultiUploadFile", map[string]string{
		"uploadFileId": uploadFileId,
		"fileMd5":      fileMd5,
		"sliceMd5":     sliceMd5,
		"lazyCheck":    "1",
		"opertype":     "3",
	}, nil)

	return err
}

// Delete 删除文件
func (d *Driver) Delete(ctx context.Context, files ...string) ([]string, error) {
	deleteFailed := make([]string, 0)
	var lastErr error

	for _, fileID := range files {
		taskInfos := []map[string]interface{}{
			{
				"fileId":   fileID,
				"fileName": "",
				"isFolder": 0,
			},
		}

		taskInfosBytes, err := json.Marshal(taskInfos)
		if err != nil {
			deleteFailed = append(deleteFailed, fileID)
			lastErr = err
			continue
		}

		form := map[string]string{
			"type":           "DELETE",
			"targetFolderId": "",
			"taskInfos":      string(taskInfosBytes),
		}

		_, err = d.request("https://cloud.189.cn/api/open/batch/createBatchTask.action", http.MethodPost, func(req *resty.Request) {
			req.SetFormData(form)
		}, nil)

		if err != nil {
			deleteFailed = append(deleteFailed, fileID)
			lastErr = err
		}
	}

	return deleteFailed, lastErr
}

// Thumb 获取缩略图
func (d *Driver) Thumb(ctx context.Context, expire *time.Time, ext string, e fs.Entity) (string, error) {
	return "", errors.New("not implemented")
}

// Source 获取下载链接
func (d *Driver) Source(ctx context.Context, e fs.Entity, args *driver.GetSourceArgs) (string, error) {
	var resp DownResp
	_, err := d.request("https://cloud.189.cn/api/portal/getFileInfo.action", http.MethodGet, func(req *resty.Request) {
		req.SetQueryParam("fileId", e.Source())
	}, &resp)
	if err != nil {
		return "", err
	}

	// 处理重定向获取最终下载链接
	client := resty.New().SetRedirectPolicy(
		resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}))

	res, err := client.R().SetHeader("User-Agent", "Mozilla/5.0").Get("https:" + resp.FileDownloadUrl)
	if err != nil {
		return "", err
	}

	downloadURL := resp.FileDownloadUrl
	if res.StatusCode() == 302 {
		downloadURL = res.Header().Get("location")
	}

	downloadURL = strings.Replace(downloadURL, "http://", "https://", 1)
	return downloadURL, nil
}

// Token 获取上传凭证
func (d *Driver) Token(ctx context.Context, uploadSession *fs.UploadSession, file *fs.UploadRequest) (*fs.UploadCredential, error) {
	return nil, errors.New("not implemented")
}

// CancelToken 取消上传凭证
func (d *Driver) CancelToken(ctx context.Context, uploadSession *fs.UploadSession) error {
	return nil
}

// CompleteUpload 完成上传
func (d *Driver) CompleteUpload(ctx context.Context, session *fs.UploadSession) error {
	return nil
}

// Capabilities 返回驱动能力
func (d *Driver) Capabilities() *driver.Capabilities {
	return capabilities
}

// MediaMeta 获取媒体元数据
func (d *Driver) MediaMeta(ctx context.Context, path, ext, language string) ([]driver.MediaMeta, error) {
	return nil, errors.New("not implemented")
}

// getSessionKey 获取会话密钥
func (d *Driver) getSessionKey() (string, error) {
	resp, err := d.request("https://cloud.189.cn/v2/getUserBriefInfo.action", http.MethodGet, nil, nil)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	sessionKey, ok := result["sessionKey"].(string)
	if !ok {
		return "", errors.New("sessionKey not found")
	}

	return sessionKey, nil
}

// getResKey 获取RSA密钥
func (d *Driver) getResKey() (string, string, error) {
	now := time.Now().UnixMilli()
	if d.rsa.Expire > now {
		return d.rsa.PubKey, d.rsa.PkId, nil
	}

	resp, err := d.request("https://cloud.189.cn/api/security/generateRsaKey.action", http.MethodGet, nil, nil)
	if err != nil {
		return "", "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", "", err
	}

	pubKey, _ := result["pubKey"].(string)
	pkId, _ := result["pkId"].(string)
	expire, _ := result["expire"].(float64)

	d.rsa.PubKey = pubKey
	d.rsa.PkId = pkId
	d.rsa.Expire = int64(expire)

	return pubKey, pkId, nil
}

// uploadRequest 上传请求
func (d *Driver) uploadRequest(uri string, form map[string]string, resp interface{}) ([]byte, error) {
	c := strconv.FormatInt(time.Now().UnixMilli(), 10)
	r := Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	l := Random("xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx")
	l = l[0 : 16+int(16*rand.Float32())]

	e := qs(form)
	data := AesEncrypt([]byte(e), []byte(l[0:16]))
	h := hex.EncodeToString(data)

	sessionKey := d.sessionKey
	signature := hmacSha1(fmt.Sprintf("SessionKey=%s&Operate=GET&RequestURI=%s&Date=%s&params=%s", sessionKey, uri, c, h), l)

	pubKey, pkId, err := d.getResKey()
	if err != nil {
		return nil, err
	}

	b := RsaEncode([]byte(l), pubKey, false)
	req := d.client.R().SetHeaders(map[string]string{
		"accept":         "application/json;charset=UTF-8",
		"SessionKey":     sessionKey,
		"Signature":      signature,
		"X-Request-Date": c,
		"X-Request-ID":   r,
		"EncryptionText": b,
		"PkId":           pkId,
	})

	if resp != nil {
		req.SetResult(resp)
	}

	res, err := req.Get("https://upload.cloud.189.cn" + uri + "?params=" + h)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(res.Body(), &result); err == nil {
		if code, ok := result["code"].(string); ok && code != "SUCCESS" {
			if msg, ok := result["msg"].(string); ok {
				return nil, errors.New(uri + "---" + msg)
			}
		}
	}

	return res.Body(), nil
}
