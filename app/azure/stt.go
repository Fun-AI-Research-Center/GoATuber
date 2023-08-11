package azure

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/proxy"
	jsoniter "github.com/json-iterator/go"
)

//azure speech-to-text服务

func GetSpeech(e *engine.Engine, speechData *bytes.Buffer) (string, error) {
	config := e.Config.Application.Azure.AzureSTT

	url := "https://" + config.SpeechRegion + ".stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=" + config.Language + "&format=" + config.Format
	req, _ := http.NewRequest("POST", url, speechData)

	req.Header.Set("Ocp-Apim-Subscription-Key", config.ApiKey)
	req.Header.Set("Content-Type", "audio/wav")

	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("获取代理客户端错误(azure-stt模块)："+er.Error()), err.Normal)
	}
	resp, er := client.Do(req)
	if er != nil {
		return "", errors.New("请求错误(azure-stt模块)：" + er.Error())
	}

	defer resp.Body.Close()
	body, er := io.ReadAll(resp.Body)
	if er != nil {
		return "", errors.New("读取返回数据错误(azure-stt模块)：" + er.Error())
	}

	var response sttResponse
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	er = json.Unmarshal(body, &response)
	if er != nil {
		return "", errors.New("json反序列化失败(azure-stt模块)：" + er.Error())
	}

	if response.RecognitionStatus != "Success" {
		return "", errors.New("azure语音识别失败，错误信息：" + response.RecognitionStatus)
	}

	return response.DisplayText, nil
}
