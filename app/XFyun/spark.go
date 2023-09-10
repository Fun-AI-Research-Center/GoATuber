package XFyun

import (
	"errors"
	"log"
	"strings"
	"time"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

func GetMessage(e *engine.Engine, message engine.PriorityMessage) error {
	if !isRead {
		getRole(e)
	}
	return generateMessage(e, message)
}

func generateMessage(e *engine.Engine, inputMessage engine.PriorityMessage) error {
	config := e.Config.Application.XunFei.XunFeiSpark

	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	URL, er := assembleAuthUrl1(config.HostURL, config.ApiKey, config.ApiSecret)
	if er != nil {
		return er
	}

	//握手并建立websocket 连接
	conn, resp, er := d.Dial(URL, nil)
	if er != nil {
		return er
	} else if resp.StatusCode != 101 {
		return errors.New("响应码错误（XFyun-spark模块）:" + resp.Status)
	}

	go func() {
		data := genParams(inputMessage.Message, inputMessage.UUID, e)
		er := conn.WriteJSON(data)
		if er != nil {
			err.Error(errors.New("发送请求错误（XFyun-spark模块）:"+er.Error()), err.Normal)
			return
		}
	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, er := conn.ReadMessage()
		if er != nil {
			return er
		}

		var (
			data map[string]interface{}
			json = jsoniter.ConfigCompatibleWithStandardLibrary
		)

		er = json.Unmarshal(msg, &data)
		if er != nil {
			return er
		}

		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			return errors.New("响应码错误（XFyun-spark模块）:" + string(msg))
		}
		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			log.Printf("Spark模型生成内容:%s\nTotalTokens:%f", answer, totalTokens)
			if int(totalTokens) > config.MemoryAndClean.ShortMemoryTokenLimit {
				go cleanMemory(e)
			}
			conn.Close()
			break
		}
	}
	//将新生成的消息压入短期记忆
	var repMessage = message{
		Role:    "assistant",
		Content: answer,
	}
	shortMemoryMessage = append(shortMemoryMessage, repMessage)

	//TODO:长期记忆

	e.Message.Message = strings.Replace(answer, "\n\n", "", 1)

	e.Ch.LLMProcess <- struct{}{}

	return nil
}

func genParams(question, uid string, e *engine.Engine) map[string]interface{} {
	config := e.Config.Application.XunFei.XunFeiSpark

	reqMessage := message{
		Role:    "user",
		Content: question,
	}
	text := append(shortMemoryMessage, reqMessage)
	text = append(roleMessage, text...)

	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": config.AppID,
			"uid":    uid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      config.Domain,
				"temperature": config.Temperature,
				"top_k":       config.TopK,
				"max_tokens":  config.MaxTokens,
				"auditing":    config.Auditing,
				"chat_id":     uid,
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": text,
			},
		},
	}
	return data
}
