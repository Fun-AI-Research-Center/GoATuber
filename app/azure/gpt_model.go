package azure

import "GoATuber-2.0/engine"

var shortMemoryMessage = make([]requestMessages, 0) //短期记忆消息
var (
	roleMessage = make([]requestMessages, 0) //角色消息
	isRead      = false                      //是否已经读取过角色信息
)

//azure openai相关model

type requestMessages struct {
	Role    string `json:"role"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type receiveMessages struct {
	Role         string `json:"role"`
	Content      string `json:"content"`
	FunctionCall struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function_call"`
}

type postDataWithFunction struct {
	Model            string            `json:"model"`
	Messages         []requestMessages `json:"messages"`
	MaxTokens        int               `json:"max_tokens"`
	Temperature      float64           `json:"temperature"`
	TopP             float64           `json:"top_p"`
	Stop             string            `json:"stop"`
	PresencePenalty  float64           `json:"presence_penalty"`
	FrequencyPenalty float64           `json:"frequency_penalty"`
	User             string            `json:"user"`
	Function         []interface{}     `json:"functions"`
	FunctionCall     string            `json:"function_call"`
}

type postDataWithoutFunction struct {
	Model            string            `json:"model"`
	Messages         []requestMessages `json:"messages"`
	MaxTokens        int               `json:"max_tokens"`
	Temperature      float64           `json:"temperature"`
	TopP             float64           `json:"top_p"`
	Stop             string            `json:"stop"`
	PresencePenalty  float64           `json:"presence_penalty"`
	FrequencyPenalty float64           `json:"frequency_penalty"`
	User             string            `json:"user"`
}

// OpenAiRcv Response
type openAiRcv struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      receiveMessages `json:"message"`
		FinishReason string          `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokes int `json:"completion_tokens"`
		TotalTokens     int `json:"total_tokens"`
	}
}

// 初始化request model
func (req *postDataWithFunction) initRequestModel(e *engine.Engine, msg engine.PriorityMessage) {
	config := e.Config.Application.Azure.AzureOpenai
	req.Model = config.Model
	//压入消息
	request := requestMessages{
		Role:    "user",
		Content: msg.Message,
		Name:    msg.UUID,
	}
	shortMemoryMessage = append(shortMemoryMessage, request)

	req.Messages = append(roleMessage, shortMemoryMessage...)
	req.FrequencyPenalty = config.FrequencyPenalty
	req.PresencePenalty = config.PresencePenalty
	req.TopP = config.TopP
	req.Temperature = config.Temperature
	req.MaxTokens = config.MaxTokens
	req.Stop = config.Stop
	req.User = msg.Username
	req.Function = e.Tool.FunctionCall.FunctionJson
	req.FunctionCall = "auto"
}

func (req *postDataWithoutFunction) initRequestModel(e *engine.Engine, msg engine.PriorityMessage) {
	config := e.Config.Application.Azure.AzureOpenai
	req.Model = config.Model
	//压入消息
	request := requestMessages{
		Role:    "user",
		Content: msg.Message,
		Name:    msg.UUID,
	}
	shortMemoryMessage = append(shortMemoryMessage, request)

	req.Messages = append(roleMessage, shortMemoryMessage...)
	req.FrequencyPenalty = config.FrequencyPenalty
	req.PresencePenalty = config.PresencePenalty
	req.TopP = config.TopP
	req.Temperature = config.Temperature
	req.MaxTokens = config.MaxTokens
	req.Stop = config.Stop
	req.User = msg.Username
}
