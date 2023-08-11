package memory

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/proxy"
	jsoniter "github.com/json-iterator/go"
)

type piencone struct {
	Type      string //消息类型
	Namespace string //命名空间
	UserName  string //用户名
	Human     string //人类发言内容
	AI        string //AI发言内容
}

// 存储请求
type requestStoreDatabaseForPinecone struct {
	Vectors   []Vector `json:"vectors"`
	Namespace string   `json:"namespace"`
}

type Vector struct {
	ID       string         `json:"id"`
	Values   []float32      `json:"values"`
	Metadata map[string]any `json:"metadata"`
}

// 储存返回
type responseStoreDatabaseForPinecone struct {
	UpsertedCount int64 `json:"upsertedCount"`
}

// 获取请求
type requestQueryDatabaseForPinecone struct {
	Filter          map[string]any `json:"filter"`
	IncludeValues   bool           `json:"includeValues"`
	IncludeMetadata bool           `json:"includeMetadata"`
	Vector          []float32      `json:"vector"`
	TopK            int            `json:"topK"`
	Namespace       string         `json:"namespace"`
}

// 获取返回
type responseQueryDatabaseForPinecone struct {
	Matches   []Matches `json:"matches"`
	Namespace string    `json:"namespace"`
}

type Matches struct {
	Id       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (p *piencone) pineconeQuery(e *engine.Engine, vector []float32) ([]string, error) {
	config := e.Config.Application.Pinecone
	url := config.Url + "/query"

	data := requestQueryDatabaseForPinecone{
		Filter: map[string]any{
			"Type": p.Type,
		},
		IncludeValues:   false,
		IncludeMetadata: true,
		Vector:          vector,
		TopK:            1,
		Namespace:       p.Namespace,
	}

	//JSON序列化
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	payload, er := json.Marshal(data)
	if er != nil {
		return nil, errors.New("JSON序列化失败(piencone):" + er.Error())
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.ApiKey)

	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("代理错误(piencone):"+er.Error()), err.Normal)
	}

	res, er := client.Do(req)
	if er != nil {
		return nil, errors.New("请求错误(piencone):" + er.Error())
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response responseQueryDatabaseForPinecone
	er = json.Unmarshal(body, &response)
	if er != nil {
		return nil, errors.New("JSON反序列化失败(piencone)，返回内容:" + string(body))
	}

	if len(response.Matches) == 0 {
		return nil, errors.New("返回内容错误(piencone):" + string(body))
	}

	return convertMapToStringSlice(response.Matches[0].Metadata), nil
}

// 空接口转字符串切片
func convertMapToStringSlice(data map[string]interface{}) []string {
	result := make([]string, 0, len(data))
	for _, value := range data {
		strValue := fmt.Sprintf("%v", value)
		result = append(result, strValue)
	}
	return result
}

func (p *piencone) pineconeStore(e *engine.Engine, vector []float32) error {
	config := e.Config.Application.Pinecone
	url := config.Url + "/vectors/upsert"

	id, er := generateRandomString(10)
	if er != nil {
		return er
	}

	data := requestStoreDatabaseForPinecone{
		Vectors: []Vector{
			{
				ID:     id,
				Values: vector,
				Metadata: map[string]any{
					"Type":  p.Type,
					"Human": p.Human,
					"AI":    p.AI,
					"User":  p.UserName,
				},
			},
		},
		Namespace: p.Namespace,
	}

	//JSON序列化
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	payload, er := json.Marshal(data)
	if er != nil {
		return errors.New("JSON序列化失败(piencone):" + er.Error())
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.ApiKey)

	client, er := proxy.Client(e)
	if er != nil {
		err.Error(errors.New("代理错误(piencone):"+er.Error()), err.Normal)
	}

	resp, er := client.Do(req)
	if er != nil {
		return errors.New("请求错误(piencone):" + er.Error())
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var response responseStoreDatabaseForPinecone
	er = json.Unmarshal(body, &response)
	if er != nil {
		return errors.New("JSON反序列化失败(piencone)，返回内容:" + string(body))
	} else if response.UpsertedCount < 1 {
		return errors.New("pinecone upsert错误:插入向量数量错误:" + strconv.Itoa(int(response.UpsertedCount)))
	}
	return nil
}

// 随机生成10位的字符串用作id
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, er := rand.Read(b)
	if er != nil {
		return "", errors.New("随机字符串生成失败(piencone):" + er.Error())
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
