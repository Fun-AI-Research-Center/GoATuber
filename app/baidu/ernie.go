package baidu

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/proxy"
	jsoniter "github.com/json-iterator/go"
)

//文心一言调用。现在默认只能使用ERNIE-Bot-turbo
//TODO:把这东西弄完

func GetMessage(e *engine.Engine, message engine.PriorityMessage) error {
	if !isRead {
		getRole(e)
	}
	if !isInitGetToken {
		go getErnieAccessToken(e)
		isInitGetToken = <-changeIsInitGetToken
	}
	return generateMessage(e, message)
}

func generateMessage(e *engine.Engine, m engine.PriorityMessage) error {
	//获取config配置文件
	config := e.Config.Application.BaiDu.BaiduErnie

	//开始构造请求体
	reqUrl := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant?access_token=" + config.AccessToken

	var (
		requestMessage request
		reqMessage     []message
	)
	tempMessage := message{
		Role:    "user",
		Content: m.Message,
	}
	shortMemoryMessage = append(shortMemoryMessage, tempMessage)
	reqMessage = append(roleMessage, shortMemoryMessage...)
	requestMessage.Messages = reqMessage
	requestMessage.Temperature = config.Temperature
	requestMessage.TopP = config.TopP
	requestMessage.PenaltyScore = config.PenaltyScore
	requestMessage.UserID = m.UUID

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	requestMessageBytes, er := json.Marshal(requestMessage)
	if er != nil {
		return er
	}

	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(requestMessageBytes))
	req.Header.Set("Content-Type", "application/json")

	client, er := proxy.Client(e)
	if er != nil {
		return er
	}

	resp, er := client.Do(req)
	if er != nil {
		return er
	}
	defer resp.Body.Close()

	var responseMessage response
	responseMessageBytes, er := io.ReadAll(resp.Body)
	if er != nil {
		return er
	}
	er = json.Unmarshal(responseMessageBytes, &responseMessage)
	if er != nil {
		return er
	}
	if responseMessage.ErrorMsg != "" {
		return errors.New("文心一言生成消息失败，错误码：" + strconv.Itoa(responseMessage.ErrorCode) + "，错误信息：" + responseMessage.ErrorMsg)
	}

	e.Message.Message = strings.Replace(responseMessage.Result, "\n\n", "", 1)
	log.Printf("文心一言生成内容:%s\n ToatalTokens:%d+%d=%d", responseMessage.Result, responseMessage.Usage.PromptTokens, responseMessage.Usage.CompletionTokens, responseMessage.Usage.TotalTokens)

	//检查短期记忆是否达到token使用上限，或文心一言是否建议清空历史记忆
	if responseMessage.Usage.TotalTokens > config.MemoryAndClean.ShortMemoryTokenLimit || responseMessage.NeedClearHistory {
		cleanMemory(e)
	}

	//短期记忆形成
	messageAI := message{
		Role:    "assistant",
		Content: responseMessage.Result,
	}
	shortMemoryMessage = append(shortMemoryMessage, messageAI)

	e.Ch.LLMProcess <- struct{}{}

	return nil
}
