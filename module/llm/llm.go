package llm

import (
	"container/heap"
	"errors"
	"regexp"

	"GoATuber-2.0/app/azure"
	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

func InitLLM(e *engine.Engine) {
	go getMessage(e)
	go processingMessage(e)
}

// 从优先队列中获取消息
func getMessage(e *engine.Engine) {
	queue := e.PriorityQueue
	for {
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

		//清空消息
		e.Message.Message = ""
		e.Message.MessageSlice = nil
	}
}

// 处理消息
// 可以考虑在这里做一个分流，是直接传递给后续的模块，还是先传递给语言模型（就是直接朗读还是生成回应的区别）
func handelMessage(e *engine.Engine, message engine.PriorityMessage) {
	chooseLLMModel(e, message)
}

func chooseLLMModel(e *engine.Engine, message engine.PriorityMessage) {
	config := e.Config.LLM
	var er error
	if config.Openai {
		//TODO:你说得对，但是我所有的key都过期了
	} else if config.AzureOpenai {
		er = azure.GetMessage(e, message)
	} else if config.Other {

	} else {
		err.Error(errors.New("错误，没有任何LLM模块被开启——假如你开启了LLM模块仍然出现此报错，请在项目页面上提交issue"), err.Fatal)
	}

	//错误处理
	if er != nil {
		err.Error(er, err.Normal)
		e.Ch.StartNext <- struct{}{}
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
		e.Message.MessageSlice = append(e.Message.MessageSlice, engine.MessageSlice{Index: i, Content: v})
	}
}
