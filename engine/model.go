package engine

import (
	"sync"

	"GoATuber-2.0/config"
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
}

type message struct {
	Message      string         //消息
	MessageSlice []MessageSlice //切分的消息
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
	ChatToFilter   chan interface{}
	FilterToLLM    chan []byte
	LLMProcess     chan struct{}
	AIFilter       chan struct{}
	FinishFilter   chan bool
	LLMToEmotion   chan struct{}
	EmotionToVoice chan struct{}
	ToFront        chan struct{}
	StartNext      chan struct{}
	WsDone         chan struct{}
}

// PriorityMessage 优先队列
type PriorityMessage struct {
	MessageType int     //消息类型
	priority    float32 //优先级
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
	emptyLock sync.Mutex    //并发写安全
}

//工具模型

type tool struct {
	FunctionCall functionCall
}

// 本地函数调用
type functionCall struct {
	IsUseFunctionCall bool                             //是否使用本地函数调用
	FunctionJson      []interface{}                    //函数调用的json
	Function          map[string]func([]string) string //函数调用的函数
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
)
