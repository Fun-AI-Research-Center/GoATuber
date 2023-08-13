package engine

import (
	"container/heap"
	"errors"

	"GoATuber-2.0/err"
	jsoniter "github.com/json-iterator/go"
)

//掌握在消息流入语言模型之前的消息传输
//使用优先队列模式，将消息加权排序

// InitTransport 初始化运输
func (e *Engine) initTransport() {
	queue := new(priorityQueue)
	queue.EmptyLock = make(chan struct{})
	e.PriorityQueue = queue
	heap.Init(&queue.Queue)
	go e.getMessageFromFilter()
	go e.getMessageFromExpend()
}

// Queue 优先队列相关

// Len 获取长度
func (queue *priorityQueueHeap) Len() int {
	return len(*queue)
}

// Less 优先级比较
func (queue *priorityQueueHeap) Less(i, j int) bool {
	return (*queue)[i].Priority > (*queue)[j].Priority //优先pop出priority值大的
}

// Swap 交换
func (queue *priorityQueueHeap) Swap(i, j int) {
	(*queue)[i], (*queue)[j] = (*queue)[j], (*queue)[i]
}

func (queue *priorityQueueHeap) Push(x interface{}) {
	*queue = append(*queue, x.(PriorityMessage))
}

func (queue *priorityQueueHeap) Pop() interface{} { //pop的效果刚刚好。弹出优先级最高中的最近加入的。
	old := *queue
	n := len(old)
	x := old[n-1]
	*queue = old[0 : n-1]
	return x
}

// FilterToLLM 从过滤模块接受消息，计算权重，放入优先队列
func (e *Engine) getMessageFromFilter() {
	for {
		select {
		case message := <-e.Ch.FilterToLLM:
			go e.handleMessageFromFilter(message)
		}
	}
}

type chatStruct struct {
	Type      int     `json:"type"`      //弹幕类型
	TimeStamp int64   `json:"timeStamp"` //弹幕读取到的时间
	ChatName  string  `json:"chatName"`  //弹幕发送者名称
	Message   string  `json:"message"`   //弹幕内容
	Price     float32 `json:"price"`     //sc
	Uid       string  `json:"uid"`       //B站用户uid
	IsValid   bool    `json:"isValid"`   //是否通过过滤
}

// 进一步处理，反序列化，计算权重，加入优先队列
func (e *Engine) handleMessageFromFilter(message []byte) {
	var (
		json = jsoniter.ConfigCompatibleWithStandardLibrary
		ms   chatStruct
	)
	er := json.Unmarshal(message, &ms)
	if er != nil {
		err.Error(errors.New("消息反序列化失败（transport-filter）:"+er.Error()), err.Normal)
		return
	}

	// 计算权重
	//弹幕类型权重设为10+sc/5；礼物权重设为sc/10
	var value float32
	switch ms.Type {
	case NormalChat:
		value = 10
	case SuperChat:
		value = 10 + ms.Price/5
	case GiftChat:
		value = ms.Price / 10
		if value < 1 { //TODO:指挥权交给用户，设定礼物感谢阈值
			return //连刷一毛礼物是吧
		}
	case Subscription:
		value = ms.Price
	}
	//TODO：加一个功能，不付费就没办法进入队列。（主包行为和框架作者无关）

	//加入优先队列
	e.PriorityQueue.Mu.Lock()
	heap.Push(&e.PriorityQueue.Queue, PriorityMessage{MessageType: ms.Type, Priority: value, Message: ms.Message, Username: ms.ChatName, UUID: ms.Uid})
	e.PriorityQueue.Mu.Unlock()

	//当监控线程因为空队列阻塞时，发送信号
	if e.PriorityQueue.IsEmpty {
		if isLock := e.PriorityQueue.EmptyMu.TryLock(); !isLock {
			return
		}
		e.PriorityQueue.EmptyLock <- struct{}{}
		e.PriorityQueue.IsEmpty = false
		e.PriorityQueue.EmptyMu.Unlock()
	}
}

// ExpendToLLM 从扩展模块接受消息，计算权重，放入优先队列
func (e *Engine) getMessageFromExpend() {
	//TODO:拓展槽还没做呢。考虑：在LLM模型后面再加一条道。
}

//优先队列维护器
