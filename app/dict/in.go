package dict

import (
	"errors"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/listen"
	jsoniter "github.com/json-iterator/go"
)

//入口

//定义类型常量（敏感词过滤or情感分析）

func InDict(e *engine.Engine, message listen.ChatStruct) {
	filterDict(e, message)
}

func filterDict(e *engine.Engine, ms listen.ChatStruct) {
	filter := New()
	er := filter.LoadWordDict("config/cfg/dist/"+e.Config.Application.Dict.FilterDictName, 0)
	if er != nil {
		err.Error(errors.New("加载过滤字典失败,"+er.Error()), err.Normal)
		return
	}
	isValid, _ := filter.Trie.validate(ms.Message)
	ms.IsValid = isValid
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	msg, er := json.Marshal(ms)
	if er != nil {
		err.Error(errors.New("JSON格式化错误（dict-filter模块）:"+er.Error()), err.Normal)
		return
	}
	e.Ch.FilterToLLM <- msg
}

func FilterAI(e *engine.Engine) {
	filter := New()
	er := filter.LoadWordDict("config/cfg/dist/"+e.Config.Application.Dict.FilterDictName, 0)
	if er != nil {
		err.Error(errors.New("加载过滤字典失败（AI）："+er.Error()), err.Fatal)
	}
	for {
		select {
		case <-e.Ch.AIFilter:
			isValid, _ := filter.Trie.validate(e.Message.Message)
			e.Ch.FinishFilter <- isValid
		}
	}
}
