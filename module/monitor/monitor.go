package monitor

import (
	"log"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/listen/bilibili"
)

// InitMonitor 初始化监控模块
func InitMonitor(e *engine.Engine) {
	//根据配置文件初始化监听程序
	if e.Config.Monitor.BiliBili {
		bilibili.InitBiliBiliListen(e)
	} else {
		log.Println("注意：没有任何监听程序被开启，无法监听任何直播间内容")
	}
}
