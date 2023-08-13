package front

import (
	"GoATuber-2.0/engine"
	jsoniter "github.com/json-iterator/go"
)

type outMessage struct {
	MessageType int       `json:"message_type"` //消息类型，1表示弹幕消息，2表示语音消息
	Sum         int       `json:"sum"`          //消息总数
	Messages    []message `json:"messages"`
	VType       int       `json:"VType"` //voice格式type，1表示http，2表示base64编码,3为二进制（编码成base64）
}

type message struct {
	Index      int    `json:"index"`      //消息序号
	Voice      string `json:"voice"`      //voice本身的数据
	Act        string `json:"act"`        //live2d动作数组名称
	Movement   string `json:"movement"`   //动作，全身的
	Expression string `json:"expression"` //表情，脸部的
	VType      int    `json:"VType"`      //voice格式type，1表示http，2表示base64编码,3为二进制（编码成base64）
}

// 将消息格式化为传送给前端的格式
func formatMessage() outMessage {
	var out = outMessage{
		MessageType: e.Message.MessageType,
		Sum:         len(e.Message.MessageSlice),
		Messages:    formatMessageSlice(e.Voice.VType),
	}
	return out
}

// 整理消息切片
func formatMessageSlice(vType int) []message {
	var out []message
	for _, v := range e.Message.MessageSlice {
		out = append(out, message{
			Index:      v.Index,
			Voice:      v.Voice,
			Act:        v.Emotion.Act,
			Movement:   v.Emotion.Movement,
			Expression: v.Emotion.Expression,
			VType:      vType,
		})
	}
	return out
}

func handelFailSpeechMessage() []byte {
	var out = outMessage{
		MessageType: engine.Speech,
		Sum:         0,
		Messages:    nil,
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	outJson, _ := json.Marshal(out)
	return outJson
}
