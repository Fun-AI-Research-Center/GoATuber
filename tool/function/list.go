package function

import "GoATuber-2.0/engine"

func InitFunction(e *engine.Engine) {
	e.Tool.FunctionCall.Function = make(map[string]func(*engine.Engine, []string) string)
	e.Tool.FunctionCall.FunctionJson = make([]interface{}, 0)
	e.Tool.FunctionCall.RequestFunction = make(chan []string, 1)
	e.Tool.FunctionCall.ResponseFunction = make(chan string, 1)
	initFunctionJson(e)
	addFunction(e)
}

func addFunction(e *engine.Engine) {
	config := e.Config.Tool

	//如启用记忆模式则调用记忆函数
	if config.Memory.UseMemory {
		addFunc(e, getMemory)
	}
}
