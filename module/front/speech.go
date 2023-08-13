package front

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func getSpeech(c *gin.Context) {
	speechFile, err := c.FormFile("audioFile")
	if err != nil {
		// 处理错误
		c.JSON(500, gin.H{
			"info": err.Error(),
		})
		return
	}

	speech, err := speechFile.Open()
	if err != nil {
		// 处理错误
		c.JSON(500, gin.H{
			"info": err.Error(),
		})
		return
	}
	defer speech.Close()

	speechData, err := ioutil.ReadAll(speech)
	if err != nil {
		// 处理错误
		c.JSON(500, gin.H{
			"info": err.Error(),
		})
		return
	}

	// 将 speechData 发送到通道中
	e.Ch.GetSpeech <- speechData
	//成功返回
	c.JSON(200, gin.H{
		"info": "ok",
	})
}
