package control

import (
	"errors"
	"os"

	"GoATuber-2.0/config"
	"GoATuber-2.0/err"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

//修改cfg目录下面的文件

func modifyDict(c *gin.Context) {
	var dict config.Dict
	if er := c.ShouldBindJSON(&dict); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(dict)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		return
	}

	file, er := os.OpenFile("./config/cfg/dict.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Application.Dict = dict
	respOK(c)
}

func modifyFilter(c *gin.Context) {
	var filter config.Filter
	if er := c.ShouldBindJSON(&filter); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(filter)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/filter.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Filter = filter
	respOK(c)
}

func modifyLLM(c *gin.Context) {
	var llm config.LLM
	if er := c.ShouldBindJSON(&llm); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(llm)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/llm.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.LLM = llm
	respOK(c)
}

func modifyMonitor(c *gin.Context) {
	var monitor config.Monitor
	if er := c.ShouldBindJSON(&monitor); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(monitor)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/monitor.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Monitor = monitor
	respOK(c)
}

func modifyMood(c *gin.Context) {
	var mood config.Mood
	if er := c.ShouldBindJSON(&mood); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(mood)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/mood.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Mood = mood
	respOK(c)
}

func modifyProxy(c *gin.Context) {
	var proxy config.Proxy
	if er := c.ShouldBindJSON(&proxy); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(proxy)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/proxy.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Proxy = proxy
	respOK(c)
}

func modifySpeech(c *gin.Context) {
	var speech config.Speech
	if er := c.ShouldBindJSON(&speech); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(speech)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/speech.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Speech = speech
	respOK(c)
}

func modifyVoice(c *gin.Context) {
	var voice config.Voice
	if er := c.ShouldBindJSON(&voice); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(voice)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/voice.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Voice = voice
	respOK(c)
}

//修改cfg/listen目录下文件

func modifyBiliBili(c *gin.Context) {
	var biliBili config.BiliBili
	if er := c.ShouldBindJSON(&biliBili); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(biliBili)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/listen/bilibili.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Listen.BiliBili = biliBili
	respOK(c)
}

//修改cfg/tool目录下的文件

func modifyMemory(c *gin.Context) {
	var memory config.Memory
	if er := c.ShouldBindJSON(&memory); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(memory)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/tool/memory.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Tool.Memory = memory
	respOK(c)
}

//修改cfg/app下的文件

func modifyAzure(c *gin.Context) {
	var azure config.Azure
	if er := c.ShouldBindJSON(&azure); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(azure)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/app/azure.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Application.Azure = azure
	respOK(c)
}

func modifyBaidu(c *gin.Context) {
	var baidu config.BaiDu
	if er := c.ShouldBindJSON(&baidu); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(baidu)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/app/baidu.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Application.BaiDu = baidu
	respOK(c)
}

func modifyOpenai(c *gin.Context) {
	var openai config.Openai
	if er := c.ShouldBindJSON(&openai); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(openai)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/app/openai.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Application.Openai = openai
	respOK(c)
}

func modifyPinecone(c *gin.Context) {
	var pinecone config.Pinecone
	if er := c.ShouldBindJSON(&pinecone); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(pinecone)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/app/pinecone.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Application.Pinecone = pinecone
	respOK(c)
}

func modifyXfyun(c *gin.Context) {
	var xfyun config.XunFei
	if er := c.ShouldBindJSON(&xfyun); er != nil {
		err.Error(errors.New("绑定json失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	//将新的配置写入配置文件
	tomlFile, er := toml.Marshal(xfyun)
	if er != nil {
		err.Error(errors.New("转换toml失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	file, er := os.OpenFile("./config/cfg/app/xfyun.toml", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		err.Error(errors.New("打开配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}
	defer file.Close()

	_, er = file.Write(tomlFile)
	if er != nil {
		err.Error(errors.New("写入配置文件失败："+er.Error()), err.Normal)
		respErr(c, "error", er)
		return
	}

	e.Config.Application.XunFei = xfyun
	respOK(c)
}
