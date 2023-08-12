package baidu

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	"GoATuber-2.0/listen"
	jsoniter "github.com/json-iterator/go"
)

type data struct {
	LogID          int    `json:"log_id"`
	Conclusion     string `json:"conclusion"`
	ConclusionType int    `json:"conclusionType"`
	ErrorMsg       string `json:"error_msg"`
}

func FilterByBaidu(e *engine.Engine, message listen.ChatStruct) {
	config := e.Config.Application.BaiDu.BaiduFilter

	reqUrl := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined?access_token=" + config.AccessToken

	message.IsValid = filter(reqUrl, message.Message)
	if !message.IsValid {
		return //直接返回得了
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	msg, er := json.Marshal(message)
	if er != nil {
		err.Error(errors.New("百度过滤失败,json格式化错误:"+er.Error()), err.Normal)
		return
	}

	e.Ch.FilterToLLM <- msg
}

func FilterAIbyBaidu(e *engine.Engine) {
	for {
		<-e.Ch.AIFilter
		//生成url
		config := e.Config.Application.BaiDu.BaiduFilter
		reqUrl := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined?access_token=" + config.AccessToken

		e.Ch.FinishFilter <- filter(reqUrl, e.Message.Message)
	}
}

func filter(reqUrl, message string) bool {
	params := url.Values{}
	params.Set("text", message)

	resp, er := http.PostForm(reqUrl, params)
	if er != nil {
		err.Error(errors.New("百度过滤失败:"+er.Error()), err.Normal)
		return false
	}

	defer resp.Body.Close()
	body, er := io.ReadAll(resp.Body)
	if er != nil {
		err.Error(errors.New("百度过滤失败:"+er.Error()), err.Normal)
		return false
	}

	var (
		json     = jsoniter.ConfigCompatibleWithStandardLibrary
		response = data{}
	)
	er = json.Unmarshal(body, &response)

	if response.ErrorMsg != "" {
		err.Error(errors.New("百度过滤失败:"+response.ErrorMsg), err.Normal)
		return false
	}

	if response.ConclusionType == 1 {
		return true
	} else {
		return false
	}
}
