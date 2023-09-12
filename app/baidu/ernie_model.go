package baidu

var shortMemoryMessage = make([]message, 0) //短期记忆消息
var (
	roleMessage = make([]message, 0) //角色消息
	isRead      = false              //是否已经读取过角色信息
)
var (
	isInitGetToken       bool
	changeIsInitGetToken = make(chan bool, 1)
)

// 传递请求结构体
type request struct {
	Messages     []message `json:"messages"`    //请求内容
	Temperature  float64   `json:"temperature"` //温度
	TopP         float64   `json:"top_p"`
	PenaltyScore float64   `json:"penalty_score"`
	UserID       string    `json:"user_id"`
}

// 回应结构体
type response struct {
	ID               string `json:"id"`                 //对话ID
	Object           string `json:"object"`             //回包类型
	Created          int64  `json:"created"`            //创建时间
	SentenceID       int    `json:"sentence_id"`        //句子ID，给流式返回用的
	IsEnd            bool   `json:"is_end"`             //是否结束。还是给流式返回用的。
	IsTruncated      bool   `json:"is_truncated"`       //是否截断。
	Result           string `json:"result"`             //对话返回结果
	NeedClearHistory bool   `json:"need_clear_history"` //是否需要清除历史记录
	BanRound         int    `json:"ban_round"`          //当need_clear_history为true时，此字段会告知第几轮对话有敏感信息，如果是当前问题，ban_round=-1
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`     //问题token数
		CompletionTokens int `json:"completion_tokens"` //回答token数
		TotalTokens      int `json:"total_tokens"`      //总token数
	} `json:"usage"` //token统计

	ErrorCode int    `json:"error_code"` //错误码
	ErrorMsg  string `json:"error_msg"`  //错误信息
}

// 传递消息结构体
type message struct {
	Role    string `json:"role"`    //角色。文心一言仅支持user和assistant两种角色。
	Content string `json:"content"` //信息内容
}
