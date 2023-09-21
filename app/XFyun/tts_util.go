package XFyun

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"GoATuber-2.0/engine"
	"github.com/gorilla/websocket"
)

const wsurl = "wss://tts-api.xfyun.cn/v2/tts"

// 与讯飞建立连接
func connectXFYun(e *engine.Engine) (*websocket.Conn, error) {
	config := e.Config.Application.XunFei.XunFeiTTS

	u := assembleAuthUrl(wsurl, config.ApiKey, config.ApiSecret)
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, resp, err := dialer.Dial(u, nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 101 {
		log.Println("连接讯飞语音接口失败，错误信息：", err)
		return nil, err
	}
	return conn, nil
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl(hostURL string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hostURL)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sign := strings.Join(signString, "\n")
	//签名结果
	sha := HmacWithShaToBase64(sign, apiSecret)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))
	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callUrl := hostURL + "?" + v.Encode()
	return callUrl
}

func HmacWithShaToBase64(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
