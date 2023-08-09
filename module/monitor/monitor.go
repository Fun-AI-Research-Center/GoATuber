package monitor

import (
	"errors"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/listen/bilibili"
)

// InitMonitor 初始化监控模块
func InitMonitor(e *engine.Engine) {
	//根据配置文件初始化监听程序
	if e.Config.Monitor.BiliBili {
		bilibili.InitBiliBiliListen(e)
	} else {
		err.Error(errors.New("未开启任何监听"), err.Fatal)
	}
}
