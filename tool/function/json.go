package function

import "GoATuber-2.0/engine"

// InitFunctionJson 初始化JSON信息
func initFunctionJson(e *engine.Engine) {
	defer func() {
		if len(e.Tool.FunctionCall.FunctionJson) == 0 {
			e.Tool.FunctionCall.IsUseFunctionCall = false
		} else {
			e.Tool.FunctionCall.IsUseFunctionCall = true
		}
	}()
	//记忆函数

}

//传递给OpenAI的JSON结构体

type getMemoryJson struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string `json:"type"`
		Properties struct {
			Chat struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"Chat"`
		} `json:"properties"`
		Required []string `json:"required"`
	} `json:"parameters"`
}
