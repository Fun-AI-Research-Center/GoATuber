package bilibili

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/listen"
)

type Monitor struct {
	servers []*ChatServer
	group   sync.WaitGroup
}

func newMonitor(c *engine.Engine) *Monitor {
	m := &Monitor{
		servers: make([]*ChatServer, 0),
	}
	chatServer, e := GetChatServer(c.Config.Listen.BiliBili.RoomID)
	if e != nil {
		err.Error(errors.New("获取弹幕服务器失败："+e.Error()), err.Fatal)
		return nil
	}
	m.servers = append(m.servers, chatServer)
	return m
}

func work(chat *ChatServer, group *sync.WaitGroup, e *engine.Engine) {
	out := make(chan Message, ChanBufSize*2) //两倍缓冲
	go chat.ReceiveMsg(out)
	for {
		msg, ok := <-out
		if !ok {
			group.Done()
			return
		}
		switch m := msg.(type) {
		case *DanMuMessage:
			chat := listen.ChatStruct{
				Type:      engine.NormalChat,
				TimeStamp: m.Timestamp,
				ChatName:  m.Uname,
				Price:     0,
				Message:   m.Text,
				Uid:       strconv.FormatInt(m.Uid, 10),
				IsValid:   true,
			}
			go chat.ChatOut(e)
		case *SuperChatMessage:
			chat := listen.ChatStruct{
				Type:      engine.SuperChat,
				TimeStamp: m.Timestamp,
				ChatName:  m.Uname,
				Price:     m.Price,
				Message:   m.Text,
				Uid:       strconv.FormatInt(m.Uid, 10),
				IsValid:   true,
			}
			go chat.ChatOut(e)
		case *GiftMessage:
			chat := listen.ChatStruct{
				Type:      engine.GiftChat,
				TimeStamp: m.Timestamp,
				ChatName:  m.Uname,
				Price:     m.Price,
				Message:   "",
				Uid:       strconv.FormatInt(m.Uid, 10),
				IsValid:   true,
			}
			go chat.ChatOut(e)
		case *GuardMessage:
			chat := listen.ChatStruct{
				Type:      engine.Subscription,
				TimeStamp: m.Timestamp,
				ChatName:  m.Uname,
				Price:     m.Price,
				Message:   "",
				Uid:       strconv.FormatInt(m.Uid, 10),
				IsValid:   true,
			}
			go chat.ChatOut(e)
			//case *EntryMessage:
			//	ifInsertError(d.insertEntryMsg(*r, m))
			//case *RoomFansMessage:
			//	ifInsertError(d.insertFansMsg(*r, m))
		}
	}
}

func (m *Monitor) start(engine *engine.Engine) {
	for _, c := range m.servers {
		for i := 1; ; i++ {
			e := c.Connect()
			if e != nil {
				err.Error(errors.New("连接直播间失败："+e.Error()+"尝试重连,连接次数："+strconv.Itoa(i)), err.Normal)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
		m.group.Add(1)
		go work(c, &(m.group), engine)
		time.Sleep(time.Second)
	}
}

func (m *Monitor) Stop() {
	for _, c := range m.servers {
		c.Disconnect()
	}
	m.group.Wait()
}
