package bilibili

import (
	"errors"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

//入口

func InitBiliBiliListen(e *engine.Engine) {
	monitor := newMonitor(e)
	if monitor == nil {
		err.Error(errors.New("初始化BiliBili监听失败"), err.Fatal)
		return
	}
	monitor.start(e)
}
