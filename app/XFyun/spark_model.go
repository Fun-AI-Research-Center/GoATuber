package XFyun

var shortMemoryMessage = make([]message, 0) //短期记忆消息
var (
	roleMessage = make([]message, 0) //角色消息
	isRead      = false              //是否已经读取过角色信息
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
