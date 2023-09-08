package control

import "github.com/gin-gonic/gin"

func respOK(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"info":   "ok",
	})
}

func respErr(c *gin.Context, status string, err error) {
	c.JSON(400, gin.H{
		"status": status,
		"info":   err.Error(),
	})
}
