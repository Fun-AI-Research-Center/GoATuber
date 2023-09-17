package control

import (
	"errors"

	"GoATuber-2.0/module/expand"
	"GoATuber-2.0/tool/split"
	"github.com/gin-gonic/gin"
)

func readText(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		respErr(c, "err", errors.New("text is empty"))
		return
	}

	expand.ReadText(e, text)

	respOK(c, nil)
}

func chat(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		respErr(c, "err", errors.New("text is empty"))
		return
	}

	expand.DirectChat(e, text)

	respOK(c, nil)
}

func getSongList(c *gin.Context) {
	list, er := expand.GetSongList()
	if er != nil {
		respErr(c, "err", er)
		return
	}

	respOK(c, list)
}

func uploadSong(c *gin.Context) {
	voice := c.PostForm("voice")
	instrument := c.PostForm("instrument")
	songName := c.PostForm("songName")

	//将string转byte
	voiceByte := []byte(voice)
	instrumentByte := []byte(instrument)

	//调用函数
	er := expand.UploadSong(voiceByte, instrumentByte, songName)
	if er != nil {
		respErr(c, "err", er)
		return
	}

	respOK(c, nil)
}

func deleteSong(c *gin.Context) {
	songName := c.PostForm("songName")

	er := expand.DeleteSong(songName)
	if er != nil {
		respErr(c, "err", er)
		return
	}

	respOK(c, nil)
}

func sing(c *gin.Context) {
	songName := c.PostForm("songName")

	er := expand.Sing(e, songName)
	if er != nil {
		respErr(c, "err", er)
		return
	}

	respOK(c, nil)
}

func getSongPiece(c *gin.Context) {
	songName := c.PostForm("songName")

	requireType := c.GetHeader("Type")
	fileRange := c.GetHeader("Range")

	var (
		piece []byte
		er    error
	)
	switch requireType {
	case "voice":
		piece, er = split.GetFilePiece("./songs/voice", songName+".wav", fileRange)

	case "instrument":
		piece, er = split.GetFilePiece("./songs/instrument", songName+".wav", fileRange)

	}
	if er != nil {
		respErr(c, "err", er)
		return
	}

	respOK(c, piece)
}
