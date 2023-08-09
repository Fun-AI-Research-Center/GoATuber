package listen

import (
	"GoATuber-2.0/engine"
)

//消息出口

// ChatOut 将消息通过引擎内管道发送到过滤器
func (chat ChatStruct) ChatOut(e *engine.Engine) {
	e.Ch.ChatToFilter <- chat
}
