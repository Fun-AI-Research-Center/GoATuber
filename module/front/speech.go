package front

import (
	"errors"

	"GoATuber-2.0/err"
	"github.com/gin-gonic/gin"
)

func getSpeech(c *gin.Context) {
	speech, er := c.GetRawData()

	if er != nil {
		err.Error(errors.New("获取speech错误："+er.Error()), err.Normal)
		c.JSON(500, gin.H{
			"info": er.Error(),
		})
		return
	}

	e.Ch.GetSpeech <- speech

	//成功返回
	c.JSON(200, gin.H{
		"info": "ok",
	})
}
