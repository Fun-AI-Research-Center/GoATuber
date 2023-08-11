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
	if e.Config.Tool.Memory.UseMemory {
		getMemoryJs := getMemoryJson{
			Name:        "getMemory",
			Description: "获取关于弹幕的历史记录。如果问题可能涉及到关于你个人的历史信息可以使用。",
			Parameters: struct {
				Type       string `json:"type"`
				Properties struct {
					Chat struct {
						Type        string `json:"type"`
						Description string `json:"description"`
					} `json:"Chat"`
				} `json:"properties"`
				Required []string `json:"required"`
			}{},
		}
		getMemoryJs.Parameters.Type = "object"
		getMemoryJs.Parameters.Properties.Chat.Type = "string"
		getMemoryJs.Parameters.Properties.Chat.Description = "对user所发送的信息的提炼"
		getMemoryJs.Parameters.Required = append(getMemoryJs.Parameters.Required, "Chat")
		addFuncJson(e, getMemoryJs)
	}
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
