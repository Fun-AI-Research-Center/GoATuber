package engine

import (
	"GoATuber-2.0/config"
)

//调度控制引擎模型，用于整个程序的调度

// NewEngine 初始化引擎
func NewEngine() *Engine {
	engine := new(Engine)
	//初始化通信系统
	initCH(engine)

	//初始化运输系统
	engine.initTransport()

	engine.Config = config.GetConfig() //配置文件导入引擎
	return engine
}

func initCH(engine *Engine) {
	engine.Ch.ChatToFilter = make(chan interface{}, 100)
	engine.Ch.FilterToLLM = make(chan []byte, 100)
	engine.Ch.LLMProcess = make(chan struct{}, 1)
	engine.Ch.FinishFilter = make(chan bool, 1)
	engine.Ch.AIFilter = make(chan struct{}, 1)
	engine.Ch.LLMToEmotion = make(chan struct{}, 1)
	engine.Ch.EmotionToVoice = make(chan struct{}, 1)
	engine.Ch.ToFront = make(chan struct{}, 1)
	engine.Ch.StartNext = make(chan struct{}, 1)
	engine.Ch.WsDone = make(chan struct{}, 1)
	engine.Ch.GetSpeech = make(chan []byte, 1)
	engine.Ch.SpeechFail = make(chan struct{}, 1)
	engine.Ch.ExpendToQueue = make(chan PriorityMessage, 1)
}
