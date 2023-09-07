package control

import (
	"errors"

	"GoATuber-2.0/module/expand"
	"github.com/gin-gonic/gin"
)

func readText(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		respErr(c, "err", errors.New("text is empty"))
		return
	}

	expand.ReadText(e, text)

	respOK(c)
}

func chat(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		respErr(c, "err", errors.New("text is empty"))
		return
	}

	expand.DirectChat(e, text)

	respOK(c)
}
