package tool

import (
	"GoATuber-2.0/engine"
	"GoATuber-2.0/tool/function"
)

// Init 初始化工具包
func Init(e *engine.Engine) {
	function.InitFunction(e)
}
