package main

import (
	"os"
	"os/signal"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/module/emotion"
	"GoATuber-2.0/module/filter"
	"GoATuber-2.0/module/front"
	"GoATuber-2.0/module/llm"
	"GoATuber-2.0/module/monitor"
	"GoATuber-2.0/module/speech"
	"GoATuber-2.0/module/voice"
	"GoATuber-2.0/tool"
)

func main() {
	//调取引擎
	e := engine.NewEngine()
	//运行引擎
	run(e)
}

func run(e *engine.Engine) {
	llm.InitLLM(e)               //先初始化获取模块
	tool.Init(e)                 //再初始化工具模块
	emotion.InitListenEmotion(e) //再初始化情感模块
	voice.InitVoice(e)           //再初始化语音模块
	go filter.InitFilter(e)      //再初始化过滤模块
	go monitor.InitMonitor(e)    //再初始化监听模块

	go speech.InitSpeech(e) //再初始化语音识别模块

	go front.InitAPI(e) //初始化前端对接

	//阻塞主线程
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}
