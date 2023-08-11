package tool

import (
	"GoATuber-2.0/engine"
	"GoATuber-2.0/tool/function"
)

// Init 初始化工具包
func Init(e *engine.Engine) {
	//一些自动激活工具的初始化
	function.InitFunction(e) //gpt函数调用

	//手动开启的初始化
}
