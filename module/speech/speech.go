package speech

import (
	"bytes"
	"container/heap"
	"errors"

	"GoATuber-2.0/app/azure"
	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

func InitSpeech(e *engine.Engine) {
	go listenSpeech(e)
}

func listenSpeech(e *engine.Engine) {
	for {
		speech := <-e.Ch.GetSpeech

		speechData := bytes.NewBuffer(speech)

		speechToText(e, speechData)
	}
}

func speechToText(e *engine.Engine, speechData *bytes.Buffer) {
	config := e.Config.Speech

	var (
		text string
		er   error
	)

	//根据配置文件选择语音转文本模块
	if config.Azure {
		text, er = azure.GetSpeech(e, speechData)
	} else if config.Other {

	} else {
		err.Error(errors.New("没有启用任何语音转文本模块"), err.Normal)
	}

	if er != nil {
		err.Error(er, err.Normal)
		return
	}

	//构造优先队列消息体
	message := engine.PriorityMessage{
		Priority:    engine.MaxPriority,
		MessageType: engine.SpeechMessage,
		Message:     text,
		Username:    "user",
		UUID:        "42",
	}

	//将消息体放入优先队列
	heap.Push(&e.PriorityQueue.Queue, message)
}
