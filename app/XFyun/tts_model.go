package XFyun

const status = 2 //数据状态，固定为2

type ttsResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Audio  string `json:"audio,omitempty"`
		Status int    `json:"status,omitempty"`
		Ced    string `json:"ced,omitempty"`
	} `json:"data"`
	Sid string `json:"sid"`
}
