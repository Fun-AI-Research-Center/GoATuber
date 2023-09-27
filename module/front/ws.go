package front

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"GoATuber-2.0/err"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws(c *gin.Context) {
	conn, er := upgrade.Upgrade(c.Writer, c.Request, nil)
	if er != nil {
		err.Error(errors.New("与web建立websocket连接失败:"+er.Error()), err.Normal)
		return
	}
	go write(conn)
	go read(conn)
}

func write(conn *websocket.Conn) {
	defer func() {
		log.Println("与展示页面断开连接")
	}()

	//向前端发送连接成功信息
	er := conn.WriteMessage(websocket.TextMessage, []byte("ok"))
	if er != nil {
		err.Error(errors.New("websocket连接建立失败:"+er.Error()), err.Normal)
	}

	for {
		select {
		case <-e.Ch.ToFront:
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			// 将OutMessage结构体转换为JSON格式
			jsonMsg, er := json.Marshal(formatMessage())
			if er != nil {
				continue
			}

			// 发送JSON消息到WebSocket连接
			er = conn.WriteMessage(websocket.TextMessage, jsonMsg)
			if er != nil {
				continue
			}
		case <-e.Ch.VoiceFail: //当语音合成失败的时候，做特殊处理。
			er = conn.WriteMessage(websocket.TextMessage, handelFailSpeechMessage())
			if er != nil {
				continue
			}
		case <-e.Ch.WsDone:
			return
		}
	}
}

type receiveMessage struct {
	Code int    `json:"code"` //消息码。0代表正确返回，其它皆为错误代码。
	Info string `json:"info"`
}

func read(conn *websocket.Conn) {
	defer func() {
		er := conn.Close()
		if er != nil {
			err.Error(errors.New("websocket连接关闭失败:"+er.Error()), err.Normal)
			return
		}
		e.Ch.WsDone <- struct{}{}
	}()
	for {
		_, msg, er := conn.ReadMessage()
		if er != nil {
			err.Error(er, err.Normal)
			e.Ch.StartNext <- struct{}{} //出错了就直接开启下一轮
			continue
		}

		var (
			receiveMsg receiveMessage
			json       = jsoniter.ConfigCompatibleWithStandardLibrary
		)

		er = json.Unmarshal(msg, &receiveMsg)
		if er != nil {
			err.Error(er, err.Normal)
			e.Ch.StartNext <- struct{}{} //出错了就直接开启下一轮
			continue
		}

		switch strconv.Itoa(receiveMsg.Code) {
		case "0":
			//清空消息
			e.Message.Message = ""
			e.Message.MessageSlice = nil
			e.Message.Username = ""
			e.Message.Uuid = ""

			e.Ch.StartNext <- struct{}{}
		default:
			err.Error(errors.New("前端返回错误代码："+strconv.Itoa(receiveMsg.Code)+",错误消息内容:"+receiveMsg.Info), err.Normal)
			e.Ch.StartNext <- struct{}{} //出错了就直接开启下一轮
		}
	}
}
