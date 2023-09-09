package expand

import (
	"GoATuber-2.0/engine"
)

//读稿子

// ReadText 绕过文本合成，直接朗读文本
func ReadText(e *engine.Engine, text string) {
	var messageType int

	messageType = engine.DirectReadWithoutMood

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
