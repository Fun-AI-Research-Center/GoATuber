package memory

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

// openai embedding
// 请求体
type requestEmbeddingForAzureOpenai struct {
	Input string `json:"input"`
	User  string `json:"model"`
}

// 响应体
type responseEmbeddingForOpenai struct {
	Object string `json:"object"`
	Data   []Data `json:"data"`
	Model  string `json:"model"`
	Usage  Usage  `json:"usage"`
}

type Data struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func embeddingByOpenAI(e *engine.Engine, message, username string) ([]float32, error) {
	//TODO:我key没钱了，，，
	return nil, nil
}

func embeddingByAzure(e *engine.Engine, message, username string) ([]float32, error) {
	config := e.Config.Application.Azure.AzureOpenai
	url := config.EndPoint + "openai/deployments/" + config.AzureEmbedding.DeploymentID + "/embeddings?api-version=" + config.AzureEmbedding.ApiVersion

	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("申请代理错误（memory-azure模块）"+er.Error()), err.Normal)
	}

	var request = requestEmbeddingForAzureOpenai{
		Input: message,
		User:  username,
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, er := json.Marshal(request)
	if er != nil {
		return nil, errors.New("json序列化错误（memory-azure模块）" + er.Error())
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.ApiKey)

	resp, er := client.Do(req)
	if er != nil {
		return nil, errors.New("请求错误（memory-azure模块）" + er.Error())
	}

	defer resp.Body.Close()
	bodyText, er := io.ReadAll(resp.Body)
	if er != nil {
		return nil, errors.New("读取错误（memory-azure模块）" + er.Error())
	}

	var response responseEmbeddingForOpenai
	er = json.Unmarshal(bodyText, &response)
	if er != nil || len(response.Data) == 0 {
		return nil, errors.New("json反序列化错误（memory-azure模块），返回消息：" + string(bodyText))
	}

	return response.Data[0].Embedding, nil
}
