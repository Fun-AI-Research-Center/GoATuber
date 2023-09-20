package engine

import (
	c "context"
	"sync"

	"GoATuber-2.0/config"
	"github.com/gin-gonic/gin"
)

//结构体模型

// Engine 引擎模型
type Engine struct {
	Message       message        //消息
	Voice         voice          //语音
	Config        config.Config  //配置文件
	Tool          tool           //工具包
	Ch            ch             //信息通道
	PriorityQueue *priorityQueue //优先队列
	R             *gin.Engine    //gin框架引擎
	Context       context        //上下文
}

type message struct {
	Username     string         //用户名
	Uuid         string         //uuid
	Message      string         //消息
	MessageSlice []MessageSlice //切分的消息
	MessageType  int            //消息类型
}

type MessageSlice struct {
	Index   int     //消息索引
	Content string  //消息内容
	Emotion emotion //情感
	Voice   string  //语音
}

type voice struct {
	VType int //voice格式type，1表示http，2表示base64编码,3为二进制（编码成base64）
}

type emotion struct {
	Emo        string //情感，按照词典进行分类
	Act        string //动作组的数组名称
	Movement   string //动作，全身的
	Expression string //表情，脸部的
}

// 信息通道
type ch struct {
	//主干流程
	ChatToFilter   chan interface{} //聊天信息传递给过滤器
	FilterToLLM    chan []byte      //过滤器传递给语言模块
	LLMProcess     chan struct{}    //语言模块处理
	AIFilter       chan struct{}    //过滤AI生成信息
	FinishFilter   chan bool        //过滤器完成
	LLMToEmotion   chan struct{}    //语言模块传递给情感模块
	EmotionToVoice chan struct{}    //情感模块传递给语音模块
	ToFront        chan struct{}    //语音模块传递给前端
	StartNext      chan struct{}    //开始下一轮
	WsDone         chan struct{}    //websocket完成

	//speech相关
	GetSpeech  chan []byte   //收到来自前端的speech信息
	SpeechFail chan struct{} //语音消息处理失败，通知前端解除阻塞

	//拓展槽
	ExpendToQueue chan PriorityMessage //拓展擦传递给优先队列
}

// Context 上下文
type context struct {
	Context c.Context
	Cancel  c.CancelFunc
}

// PriorityMessage 优先队列
type PriorityMessage struct {
	MessageType int     //消息类型
	Priority    float32 //优先级
	Message     string  //消息
	Username    string  //消息发送者名称
	UUID        string  //消息发送者识别号
}
type priorityQueueHeap []PriorityMessage

type priorityQueue struct {
	Queue     priorityQueueHeap
	Mu        sync.Mutex    //并发写安全
	IsEmpty   bool          //队列为空的时候，阻塞读取进程，此变量改为true
	EmptyLock chan struct{} //队列为空的时候启用，等待一个信号
	EmptyMu   sync.Mutex    //并发写安全
}

const MaxPriority = 1145141919810 //最大优先级

//工具模型

type tool struct {
	FunctionCall functionCall
}

// 本地函数调用
type functionCall struct {
	IsUseFunctionCall bool                                      //是否使用本地函数调用
	FunctionJson      []interface{}                             //函数调用的json
	Function          map[string]func(*Engine, []string) string //函数调用的函数
	RequestFunction   chan []string
	ResponseFunction  chan string
}

//常量类型

// 消息类型
const (
	NormalChat = iota
	SuperChat
	GiftChat
	Subscription
	SpeechMessage
	DirectReadNeedMood
	DirectReadWithoutMood
	AdministratorChatMessage
)

const (
	Chat   = 1
	Speech = 2
)
