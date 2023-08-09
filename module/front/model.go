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

func getModelInfo(c *gin.Context) {
	name := getModelName()
	if name == "" {
		return
	}
	c.JSON(http.StatusOK, name)
}

func getModelName() string {
	files, er := os.ReadDir("dist/model")
	if er != nil {
		err.Error(errors.New("读取文件错误！错误信息："+er.Error()), err.Fatal)
	}
	for _, file := range files {
		fileName := file.Name()
		match, er := regexp.Match("^.*\\.model3\\.json$", []byte(fileName))
		if er != nil {
			err.Error(errors.New("live2d初始化警告：正则表达式匹配出错！"+er.Error()), err.Normal)
		}
		if match {
			return fileName
		}
	}

	err.Error(errors.New("未找到模型文件！特别说明：该项目目前只支持model3.json后缀的live2d模型文件与vrm模型文件"), err.Fatal)
	return ""
}
