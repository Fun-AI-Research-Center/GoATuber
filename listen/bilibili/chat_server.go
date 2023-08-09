package bilibili

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// getter

func (c *ChatServer) Room() Room {
	return c.room
}

func (c *ChatServer) Host() string {
	return c.host
}

func (c *ChatServer) Port() int {
	return c.port
}

// Connect 连接弹幕服务器
func (c *ChatServer) Connect() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   4 * 1024,
		WriteBufferSize:  512,
	}
	u := fmt.Sprintf("wss://%s:%d/sub", c.host, c.wssPort)
	//请求头
	h := http.Header{}
	for name, value := range reqHeader {
		h.Add(name, value)
	}
	h.Add("Origin", "https://live.bilibili.com")
	h.Add("Cache-Control", "no-cache")
	conn, _, err := dialer.Dial(u, h)
	if err != nil {
		return err
	}
	c.conn = conn
	//进房验证
	err = c.verify()
	if err != nil {
		return err
	}
	//发送心跳包
	go c.HeartBeat()
	//读取数据包并处理
	go c.handle()
	return nil
}

func (c *ChatServer) Disconnect() {
	_ = c.conn.Close()
}

// 流水线模型 handle ==> unpackMsg ==> ReceiveMsg
func (c *ChatServer) handle() {
	unpackCh := make(chan []byte, ChanBufSize)
	go c.unpackMsg(unpackCh)
	for {
		_, buf, err := c.conn.ReadMessage()
		if err != nil {
			if buf == nil {
				close(unpackCh)
				break
			} else {
				log.Println("读取websocket消息失败", err)
			}
		}
		select {
		case unpackCh <- buf:
		default:
			log.Println("读取消息 ==> 解包数据包，阻塞！")
		}
	}
}

func (c *ChatServer) unpackMsg(in <-chan []byte) {
	for {
		msg, ok := <-in
		if !ok {
			close(c.msgCh)
			return
		}
		for _, packet := range unpack(msg) {
			select {
			case c.msgCh <- packet:
			default:
				log.Println("unpackMsg ==> ReceiveMsg，阻塞！")
			}
		}
	}
}

// ReceiveMsg 解析消息,将获取到的消息写入到 out 中
func (c *ChatServer) ReceiveMsg(out chan<- Message) {
	for {
		srcMsg, ok := <-c.msgCh
		if !ok {
			close(out)
			return
		}
		msg := parseMsg(srcMsg)
		if msg != nil {
			select {
			case out <- msg:
			default:
				log.Println("ReceiveMsg ==> out，阻塞！type:", msg.MsgType())
			}
		}
	}
}

// 发送验证消息
func (c *ChatServer) verify() error {
	verifyMsg := map[string]interface{}{
		"platform": "web",
		"protover": 3,
		"uid":      0,
		"roomid":   c.room.RealRoomID,
		"type":     2,
		"key":      c.token,
	}
	body, _ := json.Marshal(verifyMsg)
	err := c.conn.WriteMessage(websocket.BinaryMessage, pack(verInt, opEnterRoom, body))
	if err != nil {
		log.Println("发送验证信息失败:", err)
		return err
	}

	//读取服务端回传的消息，判断是否成功进入直播间，如果进入失败，服务端会断开连接
	_, buf, err := c.conn.ReadMessage()
	if err != nil {
		log.Println("读取验证信息回响失败,进入失败:", err)
		return err
	}
	op, body := unpackPacket(buf)
	if op != 8 {
		return errors.New(string(body))
	}
	return nil
}

// HeartBeat 发送心跳包,周期30s（60s不发心跳包，强行断开链接）
func (c *ChatServer) HeartBeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	ping := []byte{0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33}
	err := c.conn.WriteMessage(websocket.BinaryMessage, pack(verInt, opHeartbeat, ping))
	if err != nil {
		log.Println("ping err:", err)
		return
	}
	for range ticker.C {
		err = c.conn.WriteMessage(websocket.BinaryMessage, pack(verInt, opHeartbeat, ping))
		if err != nil {
			log.Println("ping err:", err)
			return
		}
	}
}

// GetChatServer 获得弹幕服务器地址
func GetChatServer(roomID int) (*ChatServer, error) {
	client := &BiliClient{Client: &http.Client{}} //创建新的client
	room, err := client.GetRoomInfo(roomID)
	if err != nil {
		return nil, err
	}
	val := url.Values{}
	val.Add("id", strconv.Itoa(room.RealRoomID))
	val.Add("type", "0")
	u := "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?" + val.Encode()
	resp, err := client.get(u)
	if err != nil {
		return nil, err
	}
	data := resp.Get("data")
	chatServer := &ChatServer{
		room:  room,
		msgCh: make(chan []byte, 128),
	}
	chatServer.token = data.Get("token").String()
	chatServer.host = data.Get("host_list.0").Get("host").String()
	chatServer.port = int(data.Get("host_list.0").Get("port").Int())
	chatServer.wssPort = int(data.Get("host_list.0").Get("wss_port").Int())
	chatServer.wsPort = int(data.Get("host_list.0").Get("ws_port").Int())
	return chatServer, nil
}
