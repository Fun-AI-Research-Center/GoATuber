package XFyun

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

// AssembleAuthUrl1 组装鉴权URL
func assembleAuthUrl1(hostURL, apiKey, apiSecret string) (string, error) {
	ul, err := url.Parse(hostURL)
	if err != nil {
		return "", err
	}

	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//签名内容
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sign := strings.Join(signString, "\n")
	//签名结果
	signature := hmacsha256(sign, apiSecret)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", signature)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callURL := hostURL + "?" + v.Encode()
	return callURL, nil
}

func hmacsha256(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

// 当token达到设置的上限的时候，清理记忆释放token
func cleanMemory(e *engine.Engine) {
	config := e.Config.Application.Azure.AzureOpenai.MemoryAndClean

	if config.CleanAll {
		cleanAll()
	} else if config.CleanOne {
		cleanOne()
	}
}

// 清理所有消息
func cleanAll() {
	shortMemoryMessage = shortMemoryMessage[:0]
}

// 清理切片里队首的消息
func cleanOne() {
	shortMemoryMessage = shortMemoryMessage[1:]
}

// 读取角色配置
func getRole(e *engine.Engine) {
	//从engine中获取prompt
	prompt := e.Config.Application.XunFei.XunFeiSpark.Prompt
	//当prompt为空的时候直接返回即可
	if prompt == "" {
		return
	}
	//解析prompt
	promptSlice := strings.Split(prompt, "|") //以|对句子进行分割
	for _, v := range promptSlice {
		roleSlice := strings.Split(v, ":") //以:对角色和句子进行分割
		if roleSlice[0] == "system" {      //注意：spark模型里没有system角色
			err.Error(errors.New("角色配置警告（spark模块）:spark模型不能使用system角色"), err.Normal)
			continue
		}
		var role = message{
			Role:    roleSlice[0],
			Content: strings.Join(roleSlice[1:], ":"),
		}
		roleMessage = append(roleMessage, role)
	}
	isRead = true
}
