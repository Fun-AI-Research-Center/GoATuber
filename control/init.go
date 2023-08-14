package control

import "github.com/gin-gonic/gin"

// 初始化，传递配置文件信息
func getConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"config": e.Config,
	})
}
