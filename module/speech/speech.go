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
	e.PriorityQueue.Mu.Lock()
	heap.Push(&e.PriorityQueue.Queue, message)
	e.PriorityQueue.Mu.Unlock()

	//当监控线程因为空队列阻塞时，发送信号
	if e.PriorityQueue.IsEmpty {
		if isLock := e.PriorityQueue.EmptyMu.TryLock(); !isLock {
			return
		}
		e.PriorityQueue.EmptyLock <- struct{}{}
		e.PriorityQueue.IsEmpty = false
		e.PriorityQueue.EmptyMu.Unlock()
	}
}
