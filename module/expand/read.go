package expand

import (
	"errors"
	"fmt"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

//读稿子

// ReadText 绕过文本合成，直接朗读文本
//TODO：接口调用
func ReadText(e *engine.Engine) {
	fmt.Println("请输入要朗读的文本：")
	var text string
	_, er := fmt.Scanln(&text)
	if er != nil {
		err.Error(errors.New("输入失败:"+er.Error()), err.Normal)
	}

	fmt.Println("是否要经过情感分析(y/n)：")
	var isEmotion string
	_, er = fmt.Scanln(&isEmotion)
	if er != nil {
		err.Error(errors.New("输入失败:"+er.Error()), err.Normal)
	}
	var messageType int
	if isEmotion == "y" {
		messageType = engine.DirectReadNeedMood
	} else {
		messageType = engine.DirectReadWithoutMood
	}

	//构造优先队列消息体
	message := engine.PriorityMessage{
		Priority:    engine.MaxPriority,
		MessageType: messageType,
		Message:     text,
		Username:    "user",
		UUID:        "42",
	}

	//将消息发送到优先队列
	e.Ch.ExpendToQueue <- message
}
