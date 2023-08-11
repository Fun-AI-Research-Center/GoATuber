package azure

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/proxy"
)

func GetVoice(e *engine.Engine, message *engine.MessageSlice) error {
	config := e.Config.Application.Azure.AzureTTS

	req, _ := http.NewRequest("POST", config.Url, nil)

	//设置请求头
	req.Header.Add("X-Microsoft-OutputFormat", "riff-24khz-16bit-mono-pcm")
	req.Header.Add("Content-Type", "application/ssml+xml")
	req.Header.Add("Host", config.Url)
	req.Header.Add("Authorization", config.Authentication)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit")

	//构造SSML请求体
	requestBody := []byte(`<speak version="` + config.Speak.Version + `" xmlns="` + config.Speak.Xmlns + `" xmlns:mstts="` + config.Speak.XmlnsMstts + `" xml:lang="` + config.Speak.XmlLang + `">
								<voice name="` + config.Voice.Name + `" effect="` + config.Voice.Effect + `">
									<mstts:express-as style="` + getStyle(message) + `">
										<prosody rate="` + config.Voice.Rate + `" volume="` + config.Voice.Volume + `">` +
		message.Content +
		`</prosody>
									</mstts:express-as>
								</voice>
							</speak>`)

	req.ContentLength = int64(len(requestBody))
	req.Header.Set("Content-Length", strconv.FormatInt(int64(len(requestBody)), 10))
	req.Body = io.NopCloser(bytes.NewReader(requestBody))

	//发送请求
	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("申请代理错误（azure-tts模块）:"+er.Error()), err.Normal)
	}
	resp, er := client.Do(req)
	if er != nil {
		return errors.New("发送请求错误（azure-tts模块）:" + er.Error())
	}

	defer resp.Body.Close()

	//处理响应
	bodyBytes, er := io.ReadAll(resp.Body)
	if er != nil {
		return errors.New("读取响应错误（azure-tts模块）:" + er.Error())
	}
	message.Voice = base64.StdEncoding.EncodeToString(bodyBytes)
	return nil
}

func getStyle(message *engine.MessageSlice) (style string) {
	mood := message.Emotion.Emo
	if mood == "happy" {
		style = "friendly"
	} else if mood == "mad" {
		style = "angry"
	} else if mood == "sad" {
		style = "sad"
	} else if mood == "disgust" {
		style = "disgruntled"
	} else if mood == "surprise" {
		style = "excited"
	} else if mood == "fear" {
		style = "terrified"
	} else if mood == "health" {
		style = "chat"
	}
	return style
}
