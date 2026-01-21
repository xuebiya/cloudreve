package cloud189

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// newLogin 执行登录流程
func (d *Driver) newLogin() error {
	url := "https://cloud.189.cn/api/portal/loginUrl.action?redirectURL=https%3A%2F%2Fcloud.189.cn%2Fmain.action"
	res, err := d.client.R().Get(url)
	if err != nil {
		return err
	}

	// 检查是否已登录
	redirectURL := res.RawResponse.Request.URL
	if redirectURL.String() == "https://cloud.189.cn/web/main" {
		return nil
	}

	lt := redirectURL.Query().Get("lt")
	reqId := redirectURL.Query().Get("reqId")
	appId := redirectURL.Query().Get("appId")
	headers := map[string]string{
		"lt":      lt,
		"reqid":   reqId,
		"referer": redirectURL.String(),
		"origin":  "https://open.e.189.cn",
	}

	// 获取应用配置
	var appConf AppConf
	res, err = d.client.R().SetHeaders(headers).SetFormData(map[string]string{
		"version": "2.0",
		"appKey":  appId,
	}).SetResult(&appConf).Post("https://open.e.189.cn/api/logbox/oauth2/appConf.do")
	if err != nil {
		return err
	}

	d.l.Debug("189 AppConf resp body: %s", res.String())
	if appConf.Result != "0" {
		return errors.New(appConf.Msg)
	}

	// 获取加密配置
	var encryptConf EncryptConf
	res, err = d.client.R().SetHeaders(headers).SetFormData(map[string]string{
		"appId": appId,
	}).Post("https://open.e.189.cn/api/logbox/config/encryptConf.do")
	if err != nil {
		return err
	}

	err = json.Unmarshal(res.Body(), &encryptConf)
	if err != nil {
		return err
	}

	d.l.Debug("189 EncryptConf resp body: %s", res.String())
	if encryptConf.Result != 0 {
		return errors.New("get EncryptConf error:" + res.String())
	}

	// 执行登录
	loginData := map[string]string{
		"version":         "v2.0",
		"apToken":         "",
		"appKey":          appId,
		"accountType":     appConf.Data.AccountType,
		"userName":        encryptConf.Data.Pre + RsaEncode([]byte(d.username), encryptConf.Data.PubKey, true),
		"epd":             encryptConf.Data.Pre + RsaEncode([]byte(d.password), encryptConf.Data.PubKey, true),
		"captchaType":     "",
		"validateCode":    "",
		"smsValidateCode": "",
		"captchaToken":    "",
		"returnUrl":       appConf.Data.ReturnUrl,
		"mailSuffix":      appConf.Data.MailSuffix,
		"dynamicCheck":    "FALSE",
		"clientType":      strconv.Itoa(appConf.Data.ClientType),
		"cb_SaveName":     "3",
		"isOauth2":        strconv.FormatBool(appConf.Data.IsOauth2),
		"state":           "",
		"paramId":         appConf.Data.ParamId,
	}

	res, err = d.client.R().SetHeaders(headers).SetFormData(loginData).Post("https://open.e.189.cn/api/logbox/oauth2/loginSubmit.do")
	if err != nil {
		return err
	}

	d.l.Debug("189 login resp body: %s", res.String())

	var loginResult map[string]interface{}
	if err := json.Unmarshal(res.Body(), &loginResult); err != nil {
		return err
	}

	if result, ok := loginResult["result"].(float64); !ok || int(result) != 0 {
		if msg, ok := loginResult["msg"].(string); ok {
			return errors.New(msg)
		}
		return errors.New("login failed")
	}

	return nil
}

// request 发送请求的通用方法
func (d *Driver) request(url string, method string, callback func(*resty.Request), resp interface{}) ([]byte, error) {
	var e Error
	req := d.client.R().SetError(&e).
		SetHeader("Accept", "application/json;charset=UTF-8").
		SetQueryParams(map[string]string{
			"noCache": random(),
		})

	if callback != nil {
		callback(req)
	}

	if resp != nil {
		req.SetResult(resp)
	}

	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}

	if e.ErrorCode != "" {
		if e.ErrorCode == "InvalidSessionKey" {
			err = d.newLogin()
			if err != nil {
				return nil, err
			}
			return d.request(url, method, callback, resp)
		}
		return nil, errors.New(e.ErrorMsg)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(res.Body(), &result); err == nil {
		if resCode, ok := result["res_code"].(float64); ok && int(resCode) != 0 {
			if resMsg, ok := result["res_message"].(string); ok {
				return nil, errors.New(resMsg)
			}
		}
	}

	return res.Body(), nil
}
