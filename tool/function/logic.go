package function

//函数的存储工具

import (
	"reflect"
	"runtime"

	"GoATuber-2.0/engine"
)

// GetFunctionResult 根据字符串调用函数，返回字符串
func GetFunctionResult(e *engine.Engine, functionName string, parameter []string) string {
	return executeFunction(e, get(e, functionName), parameter)
}

// 根据函数名称获取函数
func get(e *engine.Engine, funcName string) func(*engine.Engine, []string) string {
	return e.Tool.FunctionCall.Function["GoATuber-2.0/tool/function."+funcName]
}

// 执行函数
func executeFunction(engine *engine.Engine, fun func(*engine.Engine, []string) string, parameter []string) string {
	return fun(engine, parameter)
}

// 添加函数
func addFunc(e *engine.Engine, fun func(*engine.Engine, []string) string) {
	funcName := getFunctionName(fun)
	e.Tool.FunctionCall.Function[funcName] = fun
}

// 添加函数的json信息
func addFuncJson(e *engine.Engine, fun interface{}) {
	e.Tool.FunctionCall.FunctionJson = append(e.Tool.FunctionCall.FunctionJson, fun)
}

// 获取函数
func getFunctionName(fn func(*engine.Engine, []string) string) string {
	pc := reflect.ValueOf(fn).Pointer()
	function := runtime.FuncForPC(pc)
	return function.Name()
}
