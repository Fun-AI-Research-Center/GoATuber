package bilibili

import (
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type Message interface {
	MsgType() string
	setCmd(cmd string)
}

func parseMsg(src []byte) Message {
	op, body := unpackPacket(src)
	if op == opHeartbeatReply {
		//fmt.Println("heartbeat reply")
		return nil
	}
	if op != opMessage {
		return nil
	}
	result := gjson.ParseBytes(body)

	var msg Message = nil
	cmd := result.Get("cmd").String()
	switch cmd {
	case CmdDanMuMSG:
		msg = parseDanMuMessage(&result)
	case CmdSuperChatMessage:
		msg = parseSuperChatMessage(&result)
	case CmdSendGift, CmdComboSend:
		msg = parseGiftMessage(&result, cmd)
	case CmdUserToastMsg:
		msg = parseGuardMessage(&result)
	default:
		return nil
	}
	if msg != nil {
		msg.setCmd(cmd)
	}
	return msg
}

type BaseMessage struct {
	Cmd       string //原始的cmd内容
	Timestamp int64  //发送的时间戳，单位秒
}

func (r *BaseMessage) MsgType() string {
	return r.Cmd
}

func (r *BaseMessage) setCmd(c string) {
	r.Cmd = c
}

// 粉丝牌信息
type medal struct {
	MedalLevel int    //粉丝牌等级
	MedalUid   int64  //粉丝牌对应的主播
	MedalName  string //粉丝牌名称
}

// 用户信息
type user struct {
	Uid   int64  //弹幕发送者uid
	Uname string //弹幕发送者昵称
}

// DanMuMessage 弹幕消息
type DanMuMessage struct {
	BaseMessage
	medal
	user
	LiveLevel int    //弹幕发送者的直播等级
	Text      string //弹幕内容
	Types     int    //弹幕类型，滚动弹幕，底部弹幕，顶部弹幕
	FontSize  int    //字体大小
	Color     int    //弹幕颜色，10进制的rgb值
}

func parseDanMuMessage(src *gjson.Result) *DanMuMessage {
	msg := &DanMuMessage{}
	info := src.Get("info").Array()

	contentInfo := info[0].Array()
	msg.Types = int(contentInfo[1].Int())
	msg.FontSize = int(contentInfo[2].Int())
	msg.Color = int(contentInfo[3].Int())
	msg.Timestamp = contentInfo[4].Int() / 1000

	medalInfo := info[3].Array()
	if len(medalInfo) == 0 {
		//无粉丝牌信息
		msg.medal = medal{}
	} else {
		msg.MedalLevel = int(medalInfo[0].Int())
		msg.MedalName = medalInfo[1].String()
		msg.MedalUid = medalInfo[12].Int()
	}

	msg.Text = info[1].String()
	userInfo := info[2].Array()
	msg.Uid = userInfo[0].Int()
	msg.Uname = userInfo[1].String()
	msg.LiveLevel = int(src.Get("info.4.0").Int())
	return msg
}

// SuperChatMessage sc消息
type SuperChatMessage struct {
	BaseMessage
	medal
	user
	LiveLevel int     //sc发送者的直播等级
	Text      string  //sc内容
	Price     float32 //sc价格
}

func parseSuperChatMessage(src *gjson.Result) *SuperChatMessage {
	sc := &SuperChatMessage{}

	data := src.Get("data")
	sc.Timestamp = data.Get("start_time").Int()

	medalInfo := data.Get("medal_info")
	if medalInfo.Exists() {
		sc.MedalLevel = int(medalInfo.Get("medal_level").Int())
		sc.MedalName = medalInfo.Get("medal_name").String()
		sc.MedalUid = medalInfo.Get("target_id").Int()
	}

	userInfo := data.Get("user_info")
	sc.Uid = data.Get("uid").Int()
	sc.Uname = userInfo.Get("uname").String()
	sc.LiveLevel = int(userInfo.Get("user_level").Int())

	sc.Text = data.Get("message").String()
	sc.Price = float32(data.Get("price").Float())
	return sc
}

// GiftMessage 礼物消息
type GiftMessage struct {
	BaseMessage
	medal
	user
	GiftId   int     //礼物id
	GiftName string  //礼物名称
	Price    float32 //礼物价格，如果是连击则是总价值
	Num      int     //数量
}

func parseGiftMessage(src *gjson.Result, cmd string) *GiftMessage {
	gm := &GiftMessage{}
	data := src.Get("data")
	isCombo := strings.Compare(cmd, CmdComboSend) == 0
	if isCombo {
		gm.Timestamp = time.Now().Unix() //连击礼物消息中不含有时间戳信息，用当前时间代替
	} else {
		gm.Timestamp = data.Get("timestamp").Int()
	}

	medalInfo := data.Get("medal_info")
	if medalInfo.Exists() {
		gm.MedalLevel = int(medalInfo.Get("medal_level").Int())
		gm.MedalName = medalInfo.Get("medal_name").String()
		gm.MedalUid = medalInfo.Get("target_id").Int() //为0代表没牌子或者没展示
	}

	gm.Uid = data.Get("uid").Int()
	gm.Uname = data.Get("uname").String()

	if isCombo {
		gm.GiftId = int(data.Get("gift_id").Int())
		gm.GiftName = data.Get("gift_name").String()
		gm.Price = float32(data.Get("combo_total_coin").Float()) / 1000.0
		gm.Num = int(data.Get("total_num").Int())
	} else {
		gm.GiftId = int(data.Get("giftId").Int())
		gm.GiftName = data.Get("giftName").String()
		gm.Price = float32(data.Get("price").Float()) / 1000.0
		gm.Num = int(data.Get("num").Int())
	}
	return gm
}

// GuardMessage 舰长消息
type GuardMessage struct {
	BaseMessage
	user
	Name  string  //舰长，提督，总督
	Price float32 //价格
}

func parseGuardMessage(src *gjson.Result) *GuardMessage {
	gm := &GuardMessage{}

	data := src.Get("data")
	gm.Timestamp = data.Get("start_time").Int()

	gm.Uid = data.Get("uid").Int()
	gm.Uname = data.Get("username").String()
	gm.Name = data.Get("role_name").String()
	gm.Price = float32(data.Get("price").Float()) / 1000.0
	return gm
}

// EntryMessage 进场消息
type EntryMessage struct {
	BaseMessage
	user
	medal
}

// RoomFansMessage 粉丝数，粉丝团变化消息
type RoomFansMessage struct {
	BaseMessage
	Fans     int //粉丝数
	FansClub int //粉丝团
}
