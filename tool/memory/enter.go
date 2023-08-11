package memory

import (
	"errors"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

// GetMemory 获取记忆入口
func GetMemory(e *engine.Engine, message []string) string {
	config := e.Config.Tool.Memory

	var (
		embedding []float32
		memory    []string
		er        error
	)

	//embedding
	if config.UseOpenAI {

	} else if config.UseAzure {
		embedding, er = embeddingByAzure(e, message[0], message[1]) //第0号为消息，第一号为用户名
	} else {
		err.Error(errors.New("no memory embedding service is enabled"), err.Normal)
		return ""
	}

	//embedding错误处理
	if er != nil {
		err.Error(er, err.Normal)
		return ""
	}

	//search
	if config.UsePinecone {
		var p = piencone{ //TODO:更完善的结构
			Type:      "chat",
			Namespace: "live",
			UserName:  message[1],
		}
		memory, er = p.pineconeQuery(e, embedding)

	} else {
		err.Error(errors.New("no memory search service is enabled"), err.Normal)
		return ""
	}

	//search错误处理
	if er != nil {
		err.Error(er, err.Normal)
		return ""
	}

	//返回记忆
	if config.UsePinecone { //切片处理
		return "user:" + memory[1] + "说：" + memory[3] + "assistant:" + memory[0]
	}

	return memory[0] //其他默认返回第一个
}

func StoreMemory(e *engine.Engine, messageHuman, messageAI, username string) {
	config := e.Config.Tool.Memory
	var (
		embedding []float32
		er        error
	)

	//embedding
	if config.UseOpenAI {

	} else if config.UseAzure {
		embedding, er = embeddingByAzure(e, messageHuman, username)
	} else {
		err.Error(errors.New("no memory embedding service is enabled"), err.Normal)
		return
	}

	//embedding错误处理
	if er != nil {
		err.Error(er, err.Normal)
		return
	}

	//store
	if config.UsePinecone {
		var p = piencone{
			Type:      "chat",
			Namespace: "live",
			UserName:  username,
			Human:     messageHuman,
			AI:        messageAI,
		}
		er = p.pineconeStore(e, embedding)
	} else {
		err.Error(errors.New("no memory store service is enabled"), err.Normal)
		return
	}

	//store错误处理
	if er != nil {
		err.Error(er, err.Normal)
		return
	}
}
