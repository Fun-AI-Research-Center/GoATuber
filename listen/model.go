package listen

type ChatStruct struct {
	Type      int     `json:"type"`      //弹幕类型
	TimeStamp int64   `json:"timeStamp"` //弹幕读取到的时间
	ChatName  string  `json:"chatName"`  //弹幕发送者名称
	Message   string  `json:"message"`   //弹幕内容
	Price     float32 `json:"price"`     //sc
	Uid       string  `json:"uid"`       //B站用户uid
	IsValid   bool    `json:"isValid"`   //是否通过过滤
}
