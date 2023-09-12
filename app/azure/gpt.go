package azure

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/proxy"
	"GoATuber-2.0/tool/function"
	"GoATuber-2.0/tool/memory"
	jsoniter "github.com/json-iterator/go"
)

// GetMessage 从引擎中获取消息
func GetMessage(e *engine.Engine, message engine.PriorityMessage) error {
	if !isRead {
		getRole(e)
	}
	return generateMessage(e, message)
}

func generateMessage(e *engine.Engine, message engine.PriorityMessage) error {
	//初始化操作
	config := e.Config.Application.Azure.AzureOpenai
	url := config.EndPoint + "openai/deployments/" + config.DeploymentId + "/chat/completions?api-version=" + config.ApiVersion
	var postData interface{}

	//判断是否使用函数调用，初始化postdata
	if e.Tool.FunctionCall.IsUseFunctionCall {
		var postDataTemp postDataWithFunction
		postDataTemp.initRequestModel(e, message)
		postData = postDataTemp
	} else {
		var postDataTemp postDataWithoutFunction
		postDataTemp.initRequestModel(e, message)
		postData = postDataTemp
	}

	//构造请求体
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	postDataBytes, er := json.Marshal(postData)
	if er != nil {
		return errors.New("JSON格式化错误（azure-gpt模块）:" + er.Error())
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(postDataBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.ApiKey)

	//发送请求体
	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("申请代理错误（azure-gpt模块）:"+er.Error()), err.Normal)
	}
	resp, er := client.Do(req)
	if er != nil {
		return errors.New("发送请求错误（azure-gpt模块）:" + er.Error())
	}

	defer resp.Body.Close()

	if resp == nil {
		return errors.New("响应体为空（azure-gpt模块）")
	}
	body, _ := io.ReadAll(resp.Body)

	var openAiRcv openAiRcv
	er = json.Unmarshal(body, &openAiRcv)
	if er != nil {
		return errors.New("响应体反序列化错误（azure-gpt模块）:" + er.Error())
	}

	if len(openAiRcv.Choices) == 0 {
		return errors.New("azure-gpt调用失败，返回内容" + string(body))
	}

	//检查是否启用函数调用功能
	if openAiRcv.Choices[0].FinishReason == "function_call" {
		openAiRcv, er = openAiRcv.secondRequest(e, postData.(postDataWithFunction), url)
		if er != nil {
			return er
		}
	}

	log.Printf("azure openai生成内容：%s\nModel: %s TotalTokens: %d+%d=%d", openAiRcv.Choices[0].Message.Content, openAiRcv.Model, openAiRcv.Usage.PromptTokens, openAiRcv.Usage.CompletionTokes, openAiRcv.Usage.TotalTokens)
	e.Message.Message = strings.Replace(openAiRcv.Choices[0].Message.Content, "\n\n", "", 1)

	//检查短期记忆是否达到token使用上限
	if openAiRcv.Usage.TotalTokens > config.MemoryAndClean.ShortMemoryTokenLimit {
		cleanMemory(e)
	}
	//短期记忆形成
	messageAI := requestMessages{
		Role:    "assistant",
		Content: openAiRcv.Choices[0].Message.Content,
		Name:    "AI",
	}
	shortMemoryMessage = append(shortMemoryMessage, messageAI)

	//TODO:长期记忆
	if e.Config.Tool.Memory.UseMemory {
		go memory.StoreMemory(e, message.Message, messageAI.Content, message.Username)
	}

	//消息进一步处理
	e.Ch.LLMProcess <- struct{}{}
	return nil
}

func (firstResp openAiRcv) secondRequest(e *engine.Engine, firstRequest postDataWithFunction, url string) (openAiRcv, error) {
	config := e.Config.Application.Azure.AzureOpenai

	funcName := firstResp.Choices[0].Message.FunctionCall.Name

	// 定义一个正则表达式，用于匹配每一对双引号中间的内容
	pattern := `"(.*?)"`

	// 使用正则表达式解析 JSON 字符串中的字段值
	r := regexp.MustCompile(pattern)
	matches := r.FindAllStringSubmatch(firstResp.Choices[0].Message.FunctionCall.Arguments, -1)

	// 遍历字符串切片，并取出所有偶数位置上的元素
	var values []string
	for i, match := range matches {
		if i%2 == 1 { // 只保留偶数位置上的元素
			values = append(values, match[1])
		}
	}
	if values == nil {
		return openAiRcv{}, errors.New("函数调用参数解析失败（azure-gpt模块）")
	}

	result := function.GetFunctionResult(e, funcName, values)

	ms := requestMessages{
		Role:    "function",
		Name:    funcName,
		Content: result,
	}

	firstRequest.Messages = append(firstRequest.Messages, ms)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	postDataBytes, er := json.Marshal(firstRequest)
	if er != nil {
		return openAiRcv{}, errors.New("JSON格式化错误（azure-gpt模块-函数调用）:" + er.Error())
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(postDataBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.ApiKey)

	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("申请代理错误（azure-gpt模块-函数调用）:"+er.Error()), err.Normal)
	}
	resp, er := client.Do(req)
	if er != nil {
		return openAiRcv{}, errors.New("发送请求错误（azure-gpt模块-函数调用）:" + er.Error())
	}

	defer resp.Body.Close()

	if resp == nil {
		return openAiRcv{}, errors.New("响应体为空（azure-gpt模块-函数调用）")
	}
	body, _ := io.ReadAll(resp.Body)
	var secondResp openAiRcv
	er = json.Unmarshal(body, &secondResp)
	if er != nil {
		return openAiRcv{}, errors.New("响应体反序列化错误（azure-gpt模块-函数调用）:" + er.Error())
	}

	if len(secondResp.Choices) == 0 {
		return openAiRcv{}, errors.New("函数调用失败，返回内容“：" + string(body))
	}

	return secondResp, nil
}
