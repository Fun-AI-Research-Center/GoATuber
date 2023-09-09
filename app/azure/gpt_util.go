package azure

import (
	"strings"

	"GoATuber-2.0/engine"
)

// 当token达到设置的上限的时候，清理记忆释放token
func cleanMemory(e *engine.Engine) {
	config := e.Config.Application.Azure.AzureOpenai.MemoryAndClean

	if config.CleanAll {
		cleanAll()
	} else if config.CleanOne {
		cleanOne()
	}
}

// 清理所有消息
func cleanAll() {
	shortMemoryMessage = shortMemoryMessage[:0]
}

// 清理切片里队首的消息
func cleanOne() {
	shortMemoryMessage = shortMemoryMessage[1:]
}

// 读取角色配置
func getRole(e *engine.Engine) {
	//从engine中获取prompt
	prompt := e.Config.Application.Azure.AzureOpenai.Prompt
	//当prompt为空的时候直接返回即可
	if prompt == "" {
		return
	}
	//解析prompt
	promptSlice := strings.Split(prompt, "|") //以|对句子进行分割
	for _, v := range promptSlice {
		roleSlice := strings.Split(v, ":") //以:对角色和句子进行分割
		var role = requestMessages{
			Role:    roleSlice[0],
			Content: strings.Join(roleSlice[1:], ":"),
			Name:    "system",
		}
		roleMessage = append(roleMessage, role)
	}
	isRead = true
}
