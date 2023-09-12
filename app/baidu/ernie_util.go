package baidu

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	jsoniter "github.com/json-iterator/go"
)

func getErnieAccessToken(e *engine.Engine) {
	baiduErnieConfig := e.Config.Application.BaiDu.BaiduErnie

	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + baiduErnieConfig.APIKey + "&client_secret=" + baiduErnieConfig.SecretKey
	for {
		payload := strings.NewReader(``)
		client := &http.Client{}
		req, er := http.NewRequest("POST", url, payload)

		if er != nil {
			err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		res, er := client.Do(req)
		if er != nil {
			err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
		}
		defer res.Body.Close()

		body, er := io.ReadAll(res.Body)
		if er != nil {
			err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
		}

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		er = json.Unmarshal(body, &baiduErnieConfig)
		if er != nil {
			err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
		}
		if baiduErnieConfig.AccessToken == "" {
			err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
		}

		// 保存access token
		e.Config.Application.BaiDu.BaiduErnie.AccessToken = baiduErnieConfig.AccessToken
		e.Config.Application.BaiDu.BaiduErnie.ExpiresIn = baiduErnieConfig.ExpiresIn
		if !isInitGetToken {
			changeIsInitGetToken <- true
		}
		ernieClock(e)
	}
}

// 计时器，access token过期之前一个小时自动重申请access token
func ernieClock(e *engine.Engine) {
	timer := time.After(time.Duration(e.Config.Application.BaiDu.BaiduFilter.ExpiresIn-3600) * time.Second)
	<-timer
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
	prompt := e.Config.Application.BaiDu.BaiduErnie.Prompt
	//当prompt为空的时候直接返回即可
	if prompt == "" {
		return
	}
	//解析prompt
	promptSlice := strings.Split(prompt, "|") //以|对句子进行分割
	for _, v := range promptSlice {
		roleSlice := strings.Split(v, ":") //以:对角色和句子进行分割
		if roleSlice[0] == "system" {      //注意：ernie模型里没有system角色
			err.Error(errors.New("角色配置警告（ernie模块）:ernie模型不能使用system角色"), err.Normal)
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
