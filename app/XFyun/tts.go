package XFyun

import (
	"encoding/base64"
	"errors"

	"GoATuber-2.0/engine"
	jsoniter "github.com/json-iterator/go"
)

// GetVoice 讯飞语音接口
func GetVoice(e *engine.Engine, message *engine.MessageSlice) error {
	config := e.Config.Application.XunFei.XunFeiTTS

	conn, err := connectXFYun(e)
	if err != nil {
		return err
	}
	defer conn.Close()
	if conn == nil {
		return errors.New("连接讯飞语音接口失败")
	}
	request := map[string]interface{}{
		"common": map[string]interface{}{
			"app_id": config.AppID, //用户的APP_ID
		},
		"business": map[string]interface{}{
			"aue":    config.XunFeiVoice.Aue,    //音频编码
			"sfl":    config.XunFeiVoice.Sfl,    //流式返回
			"auf":    config.XunFeiVoice.Auf,    //音频采样率
			"vcn":    config.XunFeiVoice.Vcn,    //发音人
			"speed":  config.XunFeiVoice.Speed,  //语速
			"volume": config.XunFeiVoice.Volume, //音量
			"pitch":  config.XunFeiVoice.Pitch,  //音高
			"tte":    "UTF8",
			"bgs":    0,                      //是否有背景音，肯定没有啊，，，
			"reg":    config.XunFeiVoice.Reg, //英文发音方式
			"rdn":    config.XunFeiVoice.Rdn, //数字发音方式
		},
		"data": map[string]interface{}{
			"text":     base64.StdEncoding.EncodeToString([]byte(message.Content)), //文本内容
			"encoding": "UTF8",
			"status":   status, //状态
		},
	}
	err = conn.WriteJSON(request)
	if err != nil {
		return err
	}
	var voice []string
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		var (
			respMsg ttsResp
			json    = jsoniter.ConfigCompatibleWithStandardLibrary
		)
		err = json.Unmarshal(msg, &respMsg)
		if err != nil {
			return err
		}
		voice = append(voice, respMsg.Data.Audio)
		if respMsg.Data.Status == 2 {
			break
		}
	}
	var decodedVoice []byte
	for _, base := range voice {
		decoded, err := base64.StdEncoding.DecodeString(base)
		if err != nil {
			return err
		}
		decodedVoice = append(decodedVoice, decoded...)
	}
	voiceString := base64.StdEncoding.EncodeToString(decodedVoice)
	message.Voice = voiceString

	return nil
}
