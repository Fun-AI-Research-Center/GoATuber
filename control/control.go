package control

import (
	"os"
	"os/exec"
	"syscall"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/module/emotion"
	"GoATuber-2.0/module/filter"
	"GoATuber-2.0/module/front"
	"GoATuber-2.0/module/llm"
	"GoATuber-2.0/module/monitor"
	"GoATuber-2.0/module/speech"
	"GoATuber-2.0/module/voice"
	"GoATuber-2.0/tool"
	"github.com/gin-gonic/gin"
)

var (
	e *engine.Engine
)

func InitControl() {
	e = engine.NewEngine()
	initControlRouter()
}

func run(c *gin.Context) {
	llm.InitLLM(e)               //先初始化获取模块
	tool.Init(e)                 //再初始化工具模块
	emotion.InitListenEmotion(e) //再初始化情感模块
	voice.InitVoice(e)           //再初始化语音模块
	go filter.InitFilter(e)      //再初始化过滤模块
	go monitor.InitMonitor(e)    //再初始化监听模块

	go speech.InitSpeech(e) //再初始化语音识别模块

	go front.InitAPI(e) //初始化前端对接

	//调用浏览器打开展示页面
	cmd := exec.Command("cmd", "/c", "start", "http://127.0.0.1:9000")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()

	//回复前端
	respOK(c)
}

func stop(c *gin.Context) {
	os.Exit(0)
}
