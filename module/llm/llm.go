package llm

import (
	"container/heap"
	"errors"
	"regexp"
	"strconv"

	"GoATuber-2.0/app/XFyun"
	"GoATuber-2.0/app/azure"
	"GoATuber-2.0/app/baidu"
	"GoATuber-2.0/app/openai"
	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/module/voice"
)

func InitLLM(e *engine.Engine) {
	go getMessage(e)
	go processingMessage(e)
}

// 从优先队列中获取消息
func getMessage(e *engine.Engine) {
	queue := e.PriorityQueue

	for {
		//如果优先队列已经空了，那么阻塞等待，直到有新的消息进入
		if len(queue.Queue) == 0 {
			queue.IsEmpty = true
			<-queue.EmptyLock
		}

		//获取消息
		message := heap.Pop(&queue.Queue).(engine.PriorityMessage)

		//清空消息队列，保证消息时效性
		//TODO：异步衰减（这个名词疑似是chatgpt造出来的）（但是我觉得依然不是最优解，最优解就是下面这个操作）
		queue.Queue = queue.Queue[:0]
		heap.Init(&queue.Queue)

		//处理消息
		handelMessage(e, message)

		//等待下一轮
		<-e.Ch.StartNext
	}
}

// 处理消息
// 可以考虑在这里做一个分流，是直接传递给后续的模块，还是先传递给语言模型（就是直接朗读还是生成回应的区别）
func handelMessage(e *engine.Engine, message engine.PriorityMessage) {
	//分流（是否需要经过llm）
	switch message.MessageType {
	case engine.DirectReadNeedMood, engine.DirectReadWithoutMood, engine.SongMessage:
		handelDiversionMessage(e, message)
		return
	}

	//感谢礼物等

	e.Message.MessageType = getMessageType(message)
	e.Message.Username = message.Username
	e.Message.Uuid = message.UUID

	chooseLLMModel(e, message)
}

func chooseLLMModel(e *engine.Engine, message engine.PriorityMessage) {
	config := e.Config.LLM
	var er error

	if config.Openai {
		//TODO:你说得对，但是我所有的key都过期了，所以不保证没有bug
		er = openai.GetMessage(e, message)
	} else if config.AzureOpenai {
		er = azure.GetMessage(e, message)
	} else if config.XunFeiSpark {
		er = XFyun.GetMessage(e, message)
	} else if config.BaiduErnie {
		er = baidu.GetMessage(e, message)
	} else {
		err.Error(errors.New("错误，没有任何LLM模块被开启——假如你开启了LLM模块仍然出现此报错，请在项目页面上提交issue。"), err.Fatal)
	}

	//错误处理
	if er != nil {
		err.Error(er, err.Normal)
		e.Ch.StartNext <- struct{}{}
		if e.Message.MessageType == engine.Speech {
			e.Ch.SpeechFail <- struct{}{}
		}
	}
}

func processingMessage(e *engine.Engine) {
	for {
		select {
		case <-e.Ch.LLMProcess:
			e.Ch.AIFilter <- struct{}{}
			isValid := <-e.Ch.FinishFilter
			if !isValid {
				e.Message.Message = "filter!"
			}
			splitSentence(e)
			e.Ch.LLMToEmotion <- struct{}{}
		}
	}
}

func splitSentence(e *engine.Engine) {
	message := e.Message.Message

	//使用正则表达式拆局且保留分隔符号
	split := regexp.MustCompile(`([。？！；.?!;])`) //目前只有中文和英文的分隔符
	//拆句
	splitMsg := split.Split(message, -1)
	sign := split.FindAllString(message, -1)
	//去除空字符串
	for i := 0; i < len(splitMsg); i++ {
		if splitMsg[i] == "" {
			splitMsg = append(splitMsg[:i], splitMsg[i+1:]...)
			i--
		}
	}
	//将符号附加在句子后面
	for i, v := range splitMsg {
		if i < len(sign) {
			v += sign[i]
		}
		splitMsg[i] = v
	}

	//将消息整理
	for i, v := range splitMsg {
		var message = engine.MessageSlice{
			Index:   i,
			Content: v,
		}
		message.Emotion.Emo = "health"
		e.Message.MessageSlice = append(e.Message.MessageSlice, message)
	}
}

// 获得消息类型
func getMessageType(message engine.PriorityMessage) int {
	switch message.MessageType {
	case engine.NormalChat, engine.SuperChat, engine.GiftChat, engine.Subscription, engine.AdministratorChatMessage:
		return engine.Chat
	case engine.SpeechMessage:
		return engine.Speech
	default:
		err.Error(errors.New("作者疑似有点神志不清了，去提个issue叫一下他。前后端交互message-type："+strconv.Itoa(message.MessageType)), err.Normal)
		return engine.Chat
	}
}

// 处理分流消息
func handelDiversionMessage(e *engine.Engine, message engine.PriorityMessage) {
	e.Message.Message = message.Message
	splitSentence(e)

	messageType := message.MessageType
	switch messageType {
	case engine.DirectReadNeedMood:
		e.Ch.LLMToEmotion <- struct{}{}
	case engine.DirectReadWithoutMood:
		e.Ch.EmotionToVoice <- struct{}{}
	case engine.SongMessage:
		e.Message.MessageType = engine.SongMessage
		voice.HandelSongVoice(e, message)
	}
}
