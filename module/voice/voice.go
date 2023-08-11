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

func InitVoice(e *engine.Engine) {
	go listenVoice(e)
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
				go handelVoice(e, &e.Message.MessageSlice[i], &wg)
			}
			wg.Wait()

			//将消息发送到前端
			e.Ch.ToFront <- struct{}{}
		}
	}
}

func handelVoice(e *engine.Engine, message *engine.MessageSlice, wg *sync.WaitGroup) {
	config := e.Config.Voice
	var er error
	if config.XFyun {

	} else if config.Azure {
		er = azure.GetVoice(e, message)
	} else if config.Other {

	} else {
		err.Error(errors.New("错误，没有任何语音合成模块被开启——假如你开启了语音识别模块仍然出现此报错，请在项目页面上提交issue"), err.Fatal)
	}

	if er != nil {
		err.Error(er, err.Normal)
	}
	wg.Done()
}
