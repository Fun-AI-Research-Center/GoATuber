package control

import (
	"os/exec"
	"syscall"

	"github.com/gin-gonic/gin"
)

func initControlRouter() {
	r := gin.Default()
	e.R = r

	control := r.Group("/control")
	{
		//加载静态控制器页面
		control.StaticFile("/", "./dist/control.html")
		r.Static("/assets", "./dist/assets")

		//获取配置文件信息
		control.GET("/init", getConfig)

		//配置文件改动
		config := control.Group("/config")
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
			config.POST("/azure", modifyAzure) //TODO:azure需要再加上对gpt_role的设置接口
			config.POST("/baidu", modifyBaidu)
			config.POST("/openai", modifyOpenai) //TODO:openai需要再加上对gpt_role的设置接口
			config.POST("/pinecone", modifyPinecone)
			config.POST("/xunfei", modifyXfyun)
		}

		//命令相关接口
		command := control.Group("/command")
		{
			command.GET("/read", readText)
			command.GET("/chat", chat)
			song := command.Group("/song") //唱歌相关
			{
				song.GET("/list", getSongList)
				song.POST("/upload", uploadSong)
				song.DELETE("/delete", deleteSong)
				song.GET("/sing", sing)
				song.GET("/piece", getSongPiece)
			}
		}

		//启动进程
		control.GET("/run", run)
		//终止进程（目前没用——和项目结构有关系，得研究一下怎么切开
		control.GET("/stop", stop)
	}

	//调用浏览器打开前端页面
	cmd := exec.Command("cmd", "/c", "start", "http://127.0.0.1:9000/control")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()

	r.Run(":9000")
}
