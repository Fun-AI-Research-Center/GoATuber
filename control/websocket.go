package control

import (
	"errors"
	"log"
	"net/http"

	"GoATuber-2.0/err"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//通信管道

var (
	ErrToControl chan interface{}
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
		return
	}
	go write(conn)
	go read(conn)
}

func write(conn *websocket.Conn) {
	defer func() {
		log.Println("与控制器断开连接")
	}()

	for {
		select {
		case message := <-ErrToControl:
			go func() {
				msg := []byte(message.(string))
				er := conn.WriteMessage(websocket.TextMessage, msg)
				if er != nil {
					log.Println(er)
					return
				}
			}()
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
	}()

	for {
		_, message, er := conn.ReadMessage()
		if er != nil {
			break
		}

		//TODO:看看前端要返回什么
		log.Println(message)
	}
}
