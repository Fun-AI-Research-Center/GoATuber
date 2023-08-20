package config

import (
	"errors"

	"GoATuber-2.0/err"
	"github.com/spf13/viper"
)

//工具包配置

// Memory 记忆相关
type Memory struct {
	UseMemory bool `json:"use_memory" mapstructure:"UseMemory"` //使用LLM长期记忆

	//embedding相关
	UseOpenAI bool `json:"use_openai" mapstructure:"UseOpenAI"` //使用OpenAI进行embedding
	UseAzure  bool `json:"use_azure" mapstructure:"UseAzure"`   //使用Azure进行embedding

	//记忆向量库
	UsePinecone bool `json:"use_pinecone" mapstructure:"UsePinecone"` //使用Pinecone进行记忆向量库
}

func (config *Config) initMemoryConfig() {
	viper.SetConfigName("memory.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config/cfg/tool") // 指定查找配置文件的路径
	if er := viper.ReadInConfig(); er != nil {
		err.Error(errors.New("读取memory配置文件失败："+er.Error()), err.Normal)
	}
	if er := viper.Unmarshal(&config.Tool.Memory); er != nil {
		err.Error(errors.New("解析memory配置文件失败："+er.Error()), err.Normal)
	}
}
