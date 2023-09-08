package config

import (
	"errors"

	"GoATuber-2.0/err"
	"github.com/spf13/viper"
)

// TODO:疑似有个奇怪的bug。有时候数据传给前端会缺斤少两。应该不是并发问题，这里只有一条主线程。

// Config 配置文件模型
type Config struct {
	Monitor     Monitor     `json:"monitor"`
	Filter      Filter      `json:"filter"`
	LLM         LLM         `json:"llm"`
	Voice       Voice       `json:"voice"`
	Speech      Speech      `json:"speech"`
	Mood        Mood        `json:"mood"`
	Proxy       Proxy       `json:"proxy"`
	Listen      Listen      `json:"listen"`
	Application application `json:"application"`
	Tool        tool        `json:"tool"`
}

//配置模块

func initConfig() Config {
	var cfg Config
	cfg.initMonitorConfig()
	cfg.initFilterConfig()
	cfg.initLLMConfig()
	cfg.initVoiceConfig()
	cfg.initMoodConfig()
	cfg.initProxyConfig()
	cfg.initSpeechConfig()
	cfg.initApplication()
	cfg.initTool()
	return cfg
}

// 监听模块与监听包配置
func (config *Config) initMonitorConfig() {
	viper.SetConfigName("monitor.toml") // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")         // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read Monitor config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Monitor); e != nil {
		err.Error(errors.New("unmarshal Monitor config fatal"+e.Error()), err.Fatal)
	}
	config.initListener() //初始化监听
}

// 过滤模块配置
func (config *Config) initFilterConfig() {
	viper.SetConfigName("filter.toml")  // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")         // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read Filter config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Filter); e != nil {
		err.Error(errors.New("unmarshal Filter config fatal"+e.Error()), err.Fatal)
	}
}

// LLM模块配置
func (config *Config) initLLMConfig() {
	viper.SetConfigName("llm.toml")     // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")         // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read LLM config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.LLM); e != nil {
		err.Error(errors.New("unmarshal LLM config fatal"+e.Error()), err.Fatal)
	}
}

// 语音模块配置
func (config *Config) initVoiceConfig() {
	viper.SetConfigName("voice.toml")   // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")         // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read Voice config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Voice); e != nil {
		err.Error(errors.New("unmarshal Voice config fatal"+e.Error()), err.Fatal)
	}
}

// 情感模块配置（暂时没用？）
func (config *Config) initMoodConfig() {
	viper.SetConfigName("mood.toml") // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")      // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg")
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read Mood config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Mood); e != nil {
		err.Error(errors.New("unmarshal Mood config fatal"+e.Error()), err.Fatal)
	}
}

// 对话模块配置
func (config *Config) initSpeechConfig() {
	viper.SetConfigName("speech.toml") // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg")
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read Speech config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Speech); e != nil {
		err.Error(errors.New("unmarshal Speech config fatal"+e.Error()), err.Fatal)
	}
}

// 初始化代理配置
func (config *Config) initProxyConfig() {
	viper.SetConfigName("proxy.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config/cfg")
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read Proxy config fatal"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Proxy); e != nil {
		err.Error(errors.New("unmarshal Proxy config fatal"+e.Error()), err.Fatal)
	}
}

//初始化监听、应用层面模块

// 监听模块
func (config *Config) initListener() {
	if config.Monitor.BiliBili {
		//初始化BiliBili监听
		config.initBiliBili()
	} else if config.Monitor.None {
		return
	} else {
		err.Error(errors.New("no listen module is enabled"), err.Fatal)
	}
}

// 应用层面模块
func (config *Config) initApplication() {
	//初始化应用层面模块
	config.initOpenai()
	config.initPinecone()
	config.initBaidu()
	config.initAzure()
	config.initXunFei()
	config.initDict()
}

// 初始化工具包
func (config *Config) initTool() {
	config.initMemoryConfig()
}
