package openai

import "GoATuber-2.0/engine"

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
