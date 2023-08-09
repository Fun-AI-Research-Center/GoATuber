package azure

import (
	"errors"
	"io"
	"net/http"
	"time"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/proxy"
)

// GetAuthentication 获取语音合成的token
func GetAuthentication(e *engine.Engine) {
	config := e.Config.Application.Azure.AzureTTS
	fetchTokenURL := config.EndPoint
	req, _ := http.NewRequest("POST", fetchTokenURL, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", config.ApiKey)

	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("申请代理错误（azure-tts模块）:"+er.Error()), err.Normal)
	}

	//要求十分钟更新一次，故九分钟申请一次
	ticker := time.NewTimer(9 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		resp, er := client.Do(req)
		if er != nil {
			err.Error(errors.New("发送请求错误（azure-tts模块）:"+er.Error()), err.Normal)
		}

		bodyBytes, er := io.ReadAll(resp.Body)
		if er != nil {
			err.Error(errors.New("读取响应错误（azure-tts-authorization模块）:"+er.Error()), err.Normal)
		}
		accessToken := string(bodyBytes)
		authentication := "Bearer " + accessToken
		e.Config.Application.Azure.AzureTTS.Authentication = authentication
		resp.Body.Close()
	}
}
