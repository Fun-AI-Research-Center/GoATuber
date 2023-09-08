package config

import (
	"errors"

	"GoATuber-2.0/err"

	"github.com/spf13/viper"
)

//初始化监听包

//初始化BiliBili监听

// BiliBili BiliBili监听结构体
type BiliBili struct {
	RoomID int `json:"room_id" mapstructure:"RoomID"` //房间号
}

func (config *Config) initBiliBili() {
	viper.SetConfigName("BiliBili.toml")       // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")                // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/listen") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read BiliBili config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Listen.BiliBili); e != nil {
		err.Error(errors.New("unmarshal BiliBili config fatal"+e.Error()), err.Fatal)
	}
}
