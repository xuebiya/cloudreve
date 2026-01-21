package cloud189

// LoginResp 登录响应
type LoginResp struct {
	Msg    string `json:"msg"`
	Result int    `json:"result"`
	ToUrl  string `json:"toUrl"`
}

// AppConf 应用配置
type AppConf struct {
	Data struct {
		AccountType     string `json:"accountType"`
		AgreementCheck  string `json:"agreementCheck"`
		AppKey          string `json:"appKey"`
		ClientType      int    `json:"clientType"`
		IsOauth2        bool   `json:"isOauth2"`
		LoginSort       string `json:"loginSort"`
		MailSuffix      string `json:"mailSuffix"`
		PageKey         string `json:"pageKey"`
		ParamId         string `json:"paramId"`
		RegReturnUrl    string `json:"regReturnUrl"`
		ReqId           string `json:"reqId"`
		ReturnUrl       string `json:"returnUrl"`
		ShowFeedback    string `json:"showFeedback"`
		ShowPwSaveName  string `json:"showPwSaveName"`
		ShowQrSaveName  string `json:"showQrSaveName"`
		ShowSmsSaveName string `json:"showSmsSaveName"`
		Sso             string `json:"sso"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Result string `json:"result"`
}

// EncryptConf 加密配置
type EncryptConf struct {
	Result int `json:"result"`
	Data   struct {
		UpSmsOn   string `json:"upSmsOn"`
		Pre       string `json:"pre"`
		PreDomain string `json:"preDomain"`
		PubKey    string `json:"pubKey"`
	} `json:"data"`
}

// Error 错误响应
type Error struct {
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

// File 文件信息
type File struct {
	Id         int64  `json:"id"`
	LastOpTime string `json:"lastOpTime"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Icon       struct {
		SmallUrl string `json:"smallUrl"`
	} `json:"icon"`
	Url string `json:"url"`
}

// Folder 文件夹信息
type Folder struct {
	Id         int64  `json:"id"`
	LastOpTime string `json:"lastOpTime"`
	Name       string `json:"name"`
}

// Files 文件列表响应
type Files struct {
	ResCode    int    `json:"res_code"`
	ResMessage string `json:"res_message"`
	FileListAO struct {
		Count      int      `json:"count"`
		FileList   []File   `json:"fileList"`
		FolderList []Folder `json:"folderList"`
	} `json:"fileListAO"`
}

// UploadUrlsResp 上传URL响应
type UploadUrlsResp struct {
	Code       string          `json:"code"`
	UploadUrls map[string]Part `json:"uploadUrls"`
}

// Part 分片信息
type Part struct {
	RequestURL    string `json:"requestURL"`
	RequestHeader string `json:"requestHeader"`
}

// Rsa RSA密钥信息
type Rsa struct {
	Expire int64  `json:"expire"`
	PkId   string `json:"pkId"`
	PubKey string `json:"pubKey"`
}

// DownResp 下载响应
type DownResp struct {
	ResCode         int    `json:"res_code"`
	ResMessage      string `json:"res_message"`
	FileDownloadUrl string `json:"downloadUrl"`
}

// CapacityResp 容量信息响应
type CapacityResp struct {
	ResCode           int    `json:"res_code"`
	ResMessage        string `json:"res_message"`
	Account           string `json:"account"`
	CloudCapacityInfo struct {
		FreeSize     int64 `json:"freeSize"`
		MailUsedSize int64 `json:"mail189UsedSize"`
		TotalSize    int64 `json:"totalSize"`
		UsedSize     int64 `json:"usedSize"`
	} `json:"cloudCapacityInfo"`
	FamilyCapacityInfo struct {
		FreeSize  int64 `json:"freeSize"`
		TotalSize int64 `json:"totalSize"`
		UsedSize  int64 `json:"usedSize"`
	} `json:"familyCapacityInfo"`
	TotalSize uint64 `json:"totalSize"`
}
