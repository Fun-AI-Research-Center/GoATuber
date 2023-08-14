package control

import "github.com/gin-gonic/gin"

func initControlRouter() {
	r := gin.Default()
	//获取配置文件信息
	r.GET("/init", getConfig)

	//配置文件改动
	config := r.Group("/config")
	{
		config.POST("/dict", modifyDict)
		config.POST("/filter", modifyFilter)
		config.POST("/llm", modifyLLM)
		config.POST("/monitor", modifyMonitor)
		config.POST("/mood", modifyMood)
		config.POST("/proxy", modifyProxy)
		config.POST("/speech", modifySpeech)
		config.POST("/voice", modifyVoice)
		//监听包
		config.POST("/bilibili", modifyBiliBili)
		//工具包
		config.POST("/memory", modifyMemory)
		//应用包
		config.POST("/azure", modifyAzure)
		config.POST("/baidu", modifyBaidu)
		config.POST("/openai", modifyOpenai)
		config.POST("/pinecone", modifyPinecone)
		config.POST("/xunfei", modifyXfyun)
	}

	//命令相关接口
	command := r.Group("/command")
	{
		command.GET("/read", readText)
	}

	//启动进程
	r.GET("/run", run)
	//终止进程
	r.GET("/stop", stop)
	r.Run(":8900")
}
