package front

import (
	"os"
	"strings"

	"GoATuber-2.0/engine"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var e *engine.Engine

func InitAPI(engine *engine.Engine) {
	e = engine
	r := gin.Default()
	r.GET("/ws", ws)
	r.Use(static.Serve("/", static.LocalFile("dist", true)))
	r.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := os.ReadFile("./dist/index.html")
			if err != nil {
				c.Writer.WriteHeader(404)
				c.Writer.WriteString("Not Found")
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Header().Add("Accept", "text/html")
			c.Writer.Write(content)
			c.Writer.Flush()
		}
	})
	r.GET("/get", getModelInfo)
	r.POST("/speech", getSpeech)
	r.Run(":9000") //服务在本地9000端口运行
}
