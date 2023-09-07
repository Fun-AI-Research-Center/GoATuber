package voice

import (
	"errors"
	"sync"

	"GoATuber-2.0/app/azure"
	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

// voice type
const (
	url    = 1
	base64 = 2
	binary = 3
)

var (
	errorHandel   = make(chan struct{}, 1)
	getErrorInfo  = make(chan struct{}, 1)
	postErrorInfo = make(chan bool, 1)
)

func InitVoice(e *engine.Engine) {
	go listenVoice(e)
	go handelError()
	if e.Config.Voice.Azure {
		e.Voice.VType = base64
		go azure.GetAuthentication(e)
		azure.GetTTSUrl(e)
	}
}

func listenVoice(e *engine.Engine) {
	for {
		select {
		case <-e.Ch.EmotionToVoice:
			var wg sync.WaitGroup
			for i, _ := range e.Message.MessageSlice {
				wg.Add(1)
				go handelVoice(e, &e.Message.MessageSlice[i], &wg, 0)
			}
			wg.Wait()
			getErrorInfo <- struct{}{}

			//错误信息处理
			isErr := <-postErrorInfo
			if isErr {
				e.Ch.StartNext <- struct{}{}
				if e.Message.MessageType == engine.Speech {
					e.Ch.SpeechFail <- struct{}{}
				}
				continue
			}

			//将消息发送到前端
			e.Ch.ToFront <- struct{}{}
		}
	}
}

func handelVoice(e *engine.Engine, message *engine.MessageSlice, wg *sync.WaitGroup, tryCount int) {
	config := e.Config.Voice
	var er error
	if config.XFyun {

	} else if config.Azure {
		er = azure.GetVoice(e, message)
	} else {
		err.Error(errors.New("错误，没有任何语音合成模块被开启——假如你开启了语音识别模块仍然出现此报错，请在项目页面上提交issue"), err.Fatal)
	}

	if er != nil {
		err.Error(er, err.Normal)
		//进行重试
		wg.Add(1)
		tryCount++
		if tryCount < 3 { //重试三次
			handelVoice(e, message, wg, tryCount)
		} else {
			errorHandel <- struct{}{}
			wg.Done()
		}
	}

	wg.Done()
}

// 处理错误
func handelError() {
	var isErr = false
	for {
		select {
		case <-errorHandel:
			isErr = true
		case <-getErrorInfo:
			postErrorInfo <- isErr
			isErr = false
		}
	}
}
