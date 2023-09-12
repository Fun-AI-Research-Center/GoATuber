package filter

import (
	"GoATuber-2.0/app/baidu"
	"GoATuber-2.0/app/dict"
	"GoATuber-2.0/engine"
	"GoATuber-2.0/listen"
)

func InitFilter(e *engine.Engine) {
	go filter(e)
	if e.Config.Filter.Dict {
		go dict.FilterAI(e) //这算个结构性的失败，我失算了
	}

	if e.Config.Filter.Baidu {
		go baidu.GetFilterAccessToken(e)
		go baidu.FilterAIbyBaidu(e)
	}
}

func filter(e *engine.Engine) {
	for {
		select {
		case chat := <-e.Ch.ChatToFilter:
			go handelChat(chat.(listen.ChatStruct), e)
		}
	}
}

func handelChat(chat listen.ChatStruct, e *engine.Engine) {
	if e.Config.Filter.Dict {
		dict.InDict(e, chat)
	} else if e.Config.Filter.Baidu {
		baidu.FilterByBaidu(e, chat)
	} else {
		//不过滤
	}
}
