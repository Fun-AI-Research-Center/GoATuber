package baidu

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
	jsoniter "github.com/json-iterator/go"
)

func GetAccessToken(e *engine.Engine) {
	baiduFilterCfg := e.Config.Application.BaiDu.BaiduFilter

	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + baiduFilterCfg.APIKey + "&client_secret=" + baiduFilterCfg.SecretKey
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, er := http.NewRequest("POST", url, payload)

	if er != nil {
		err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, er := client.Do(req)
	if er != nil {
		err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
	}
	defer res.Body.Close()

	body, er := io.ReadAll(res.Body)
	if er != nil {
		log.Fatalf("%v", er)
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	er = json.Unmarshal(body, &baiduFilterCfg)
	if er != nil {

	}
	if baiduFilterCfg.AccessToken == "" {
		err.Error(errors.New("获取access token失败:"+er.Error()), err.Normal)
	}

	// 保存access token
	e.Config.Application.BaiDu.BaiduFilter.AccessToken = baiduFilterCfg.AccessToken
	e.Config.Application.BaiDu.BaiduFilter.ExpiresIn = baiduFilterCfg.ExpiresIn
	go clock(e)
}

// 计时器，access token过期之前一个小时自动重申请access token
func clock(e *engine.Engine) {
	timer := time.After(time.Duration(e.Config.Application.BaiDu.BaiduFilter.ExpiresIn-3600) * time.Second)
	<-timer
	go GetAccessToken(e)
}
