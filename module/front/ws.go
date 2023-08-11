package front

import (
	"errors"
	"log"
	"net/http"

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
		log.Println("连接断开")
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

		case <-e.Ch.WsDone:
			return
		}
	}
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
		_, code, er := conn.ReadMessage()
		if er != nil {
			break
		}
		switch string(code) {
		case "0":
			//清空消息
			e.Message.Message = ""
			e.Message.MessageSlice = nil
			e.Message.Username = ""
			e.Message.Uuid = ""

			e.Ch.StartNext <- struct{}{}
		default:
			err.Error(errors.New("前端返回错误代码："+string(code)), err.Normal)
		}
	}
}
