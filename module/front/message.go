package front

type outMessage struct {
	Sum      int       `json:"sum"` //消息总数
	Messages []message `json:"messages"`
	VType    int       `json:"VType"` //voice格式type，1表示http，2表示base64编码,3为二进制（编码成base64）
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
		Sum:      len(e.Message.MessageSlice),
		Messages: formatMessageSlice(e.Voice.VType),
	}
	return out
}

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
