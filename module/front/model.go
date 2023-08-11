package front

import (
	"errors"
	"net/http"
	"os"
	"regexp"

	"GoATuber-2.0/err"
	"github.com/gin-gonic/gin"
)

//获取模型信息，并不是model层

const (
	null   = 0 //错误
	live2d = 1 //live2d模型
	vrm    = 2 //vrm模型
)

func getModelInfo(c *gin.Context) {
	name, model := getModelName(c)
	if name == "" {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"type": model,
	})
}

func getModelName(c *gin.Context) (string, int) {
	files, er := os.ReadDir("dist/model")
	if er != nil {
		err.Error(errors.New("读取文件错误！错误信息："+er.Error()), err.Fatal)
	}
	for _, file := range files {
		fileName := file.Name()

		matchVrm, er := regexp.Match("^.*\\.vrm$", []byte(fileName))
		if er != nil {
			err.Error(errors.New("vrm初始化警告：正则表达式匹配出错！"+er.Error()), err.Normal)
		}
		if matchVrm {
			return fileName, vrm
		}

		matchLive2d, er := regexp.Match("^.*\\.model3\\.json$", []byte(fileName))
		if er != nil {
			err.Error(errors.New("live2d初始化警告：正则表达式匹配出错！"+er.Error()), err.Normal)
		}
		if matchLive2d {
			return fileName, live2d
		}
	}

	//未找到模型文件
	c.JSON(http.StatusBadRequest, gin.H{
		"info": "未找到模型文件",
	})
	err.Error(errors.New("未找到模型文件！特别说明：该项目目前只支持model3.json后缀的live2d模型文件与vrm模型文件"), err.Fatal)
	return "", null
}
