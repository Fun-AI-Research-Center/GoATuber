package azure

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
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
func getRole() {
	file, er := os.Open("./config/cfg/dict/gpt_role.txt")
	if er != nil {
		err.Error(errors.New("读取角色配置文件失败:"+er.Error()), err.Normal)
		isRead = true
		return
	}
	defer file.Close()

	// 读取角色配置文件
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		texts := strings.Split(text, ":")
		var role = requestMessages{
			Role:    texts[0],
			Content: strings.Join(texts[1:], ":"),
			Name:    "system",
		}
		roleMessage = append(roleMessage, role)
	}
	isRead = true
}
