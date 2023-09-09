package expand

import "GoATuber-2.0/engine"

//直接输入文字与模型对话

func DirectChat(e *engine.Engine, text string) {
	//构造优先队列消息体
	messageType := engine.PriorityMessage{
		Priority:    engine.MaxPriority,
		MessageType: engine.AdministratorChatMessage,
		Message:     text,
		Username:    "user",
		UUID:        "42",
	}

	//将消息发送到优先队列
	e.Ch.ExpendToQueue <- messageType
}
