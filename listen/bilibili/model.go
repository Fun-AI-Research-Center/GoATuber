package bilibili

import "github.com/gorilla/websocket"

// 数据类型
const (
	verPlain  = 0 //普通文本，utf-8编码
	verInt    = 1 //控制通信消息
	verZlib   = 2 // zlib格式压缩，已弃用
	verBrotli = 3 //brotli压缩
)

// 操作码
const (
	opHeartbeat      = 2 //发送心跳包
	opHeartbeatReply = 3 //服务端回应心跳包
	opMessage        = 5 //弹幕消息等
	opEnterRoom      = 7 //认证包
	opEnterRoomReply = 8 //进入直播间成功
)

// Cmd
const (
	CmdEntryEffect      = "ENTRY_EFFECT"       //舰长进场消息
	CmdSendGift         = "SEND_GIFT"          //投喂礼物
	CmdComboSend        = "COMBO_SEND"         //礼物连击
	CmdDanMuMSG         = "DANMU_MSG"          //弹幕
	CmdUserToastMsg     = "USER_TOAST_MSG"     //续费舰长
	CmdSuperChatMessage = "SUPER_CHAT_MESSAGE" //sc
)

const (
	ChanBufSize = 64
)

// Room 直播间房间号
type Room struct {
	RoomID     int //房间号
	RealRoomID int //真的房间号
}

// ChatServer 聊天服务器
type ChatServer struct {
	room    Room //对应的直播间
	host    string
	port    int
	wssPort int
	wsPort  int
	token   string
	conn    *websocket.Conn //websocket链接
	msgCh   chan []byte     //收到的数据包，已经过解压、拆包
}
