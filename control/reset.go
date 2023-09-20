package control

import (
	"GoATuber-2.0/module/emotion"
	"GoATuber-2.0/module/filter"
	"GoATuber-2.0/module/llm"
	"GoATuber-2.0/module/monitor"
	"GoATuber-2.0/module/speech"
	"GoATuber-2.0/module/voice"
	"GoATuber-2.0/tool"
	"github.com/gin-gonic/gin"
)

// 配置文件热改动

// 主程序是否已经启动
var (
	isProgramRunning = false
	nextChange       = make(chan struct{})
)

// 中间件，检查是否需要进行热重载
func checkReset(c *gin.Context) {
	c.Next() //先执行后面的函数
	if isProgramRunning {
		<-nextChange
		reset()
	}
}

func reset() {
	//关闭所有(除了前后端对接)监听式goroutine,以防重启时出现双开
	e.Context.Cancel()

	//热重载(正确的,就是把run的命令复制了一下,非常粗暴全局重载)
	llm.InitLLM(e)               //先初始化获取模块
	tool.Init(e)                 //再初始化工具模块
	emotion.InitListenEmotion(e) //再初始化情感模块
	voice.InitVoice(e)           //再初始化语音模块
	go filter.InitFilter(e)      //再初始化过滤模块
	go monitor.InitMonitor(e)    //再初始化监听模块

	go speech.InitSpeech(e) //再初始化语音识别模块
}
