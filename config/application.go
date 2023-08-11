package config

import (
	"errors"

	"GoATuber-2.0/err"

	"github.com/spf13/viper"
)

//初始化应用配置

//综合服务（azure）

// azure结构体
type azure struct {
	AzureOpenai azureOpenai //azure openai,文本生成服务
	AzureTTS    azureTTS    //azure tts,文本转语音服务
	AzureSTT    azureSTT    //azure stt,语音转文本服务
}

// azure openai结构体
type azureOpenai struct {
	ApiKey           string  `mapstructure:"api_key"`           //API密钥
	EndPoint         string  `mapstructure:"end_point"`         //终结点，在azure openai控制台查询
	DeploymentId     string  `mapstructure:"deployment_id"`     //模型部署的名称
	ApiVersion       string  `mapstructure:"api_version"`       //要用于此操作的 API 版本
	Model            string  `mapstructure:"model"`             //模型
	Temperature      float64 `mapstructure:"temperature"`       //温度
	TopP             float64 `mapstructure:"top_p"`             //top_p
	MaxTokens        int     `mapstructure:"max_tokens"`        //最大token
	Stop             string  `mapstructure:"stop"`              //停止标志
	PresencePenalty  float64 `mapstructure:"presence_penalty"`  //presence_penalty
	FrequencyPenalty float64 `mapstructure:"frequency_penalty"` //frequency_penalty

	AzureEmbedding azureEmbedding //azure embedding,嵌入服务
}

// azure embedding azure-openai嵌入服务
type azureEmbedding struct {
	DeploymentID string `mapstructure:"deployment_id"` //模型部署的名称
	ApiVersion   string `mapstructure:"api_version"`   //要用于此操作的 API 版本
}

// azure tts结构体
type azureTTS struct {
	ApiKey   string `mapstructure:"api_key"`   //API密钥
	EndPoint string `mapstructure:"end_point"` //终结点，在azure tts控制台查询
	Speak    struct {
		Version    string `mapstructure:"version"`     //指示用于解释文档标记的 SSML 规范的版本。 当前版本为“1.0”。
		XmlLang    string `mapstructure:"xml_lang"`    //根文档的语言。 该值可以包含语言代码，例如 en（英语），也可以包含区域设置，例如 en-US（美国英语）。
		XmlnsMstts string `mapstructure:"xmlns_mstts"` //xmlns:mstts
		Xmlns      string `mapstructure:"xmlns"`       //xmlns
	}
	Voice struct {
		Name   string `mapstructure:"name"`   //用于文本转语音输出的语音。 有关受支持的预生成语音的完整列表，请参阅语言支持(https://learn.microsoft.com/zh-cn/azure/cognitive-services/speech-service/language-support?tabs=tts)。
		Effect string `mapstructure:"effect"` //音频效果处理器，用于在设备上针对特定方案优化合成语音输出的质量。 可选。
		Rate   string `mapstructure:"rate"`   //语速，填写比默认值高或低的百分比
		Volume string `mapstructure:"volume"` //音量，填写比默认值高或低的百分比
	}
	Authentication string //令牌
	Url            string //请求地址
}

// azure stt结构体
type azureSTT struct {
	ApiKey       string `mapstructure:"api_key"`       //API密钥
	SpeechRegion string `mapstructure:"speech_region"` //语音区域
	Language     string `mapstructure:"language"`      //语言
	Format       string `mapstructure:"format"`        //格式
}

// 初始化Azure配置
func (config *Config) initAzure() {
	viper.SetConfigName("azure.toml")       // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/app") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read azure config fatal:"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Application.Azure); e != nil {
		err.Error(errors.New("unmarshal azure config fatal:"+e.Error()), err.Fatal)
	}
}

//综合服务（讯飞）（考虑到今后可能接入讯飞的服务）

// 讯飞结构体
type xunFei struct {
	XunFeiTTS xunFeiTTS //讯飞 tts,文本转语音服务
}

type xunFeiTTS struct {
	ApiKey      string `mapstructure:"api_key"`    //API密钥
	ApiSecret   string `mapstructure:"api_secret"` //API密钥
	AppID       string `mapstructure:"app_id"`     //API密钥
	XunFeiVoice struct {
		Aue    string `mapstructure:"aue"`    //音频编码
		Sfl    int    `mapstructure:"sfl"`    //流式返回
		Auf    string `mapstructure:"auf"`    //音频采样率
		Vcn    string `mapstructure:"vcn"`    //发音人
		Speed  int    `mapstructure:"speed"`  //语速
		Volume int    `mapstructure:"volume"` //音量
		Pitch  int    `mapstructure:"pitch"`  //音高
		Reg    string `mapstructure:"reg"`    //英文发音方式
		Rdn    string `mapstructure:"rdn"`    //数字发音方式
	}
}

// 初始化讯飞配置
func (config *Config) initXunFei() {
	viper.SetConfigName("xfyun.toml")       // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/app") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read xunfei config fatal:"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Application.XunFei); e != nil {
		err.Error(errors.New("unmarshal xunfei config fatal:"+e.Error()), err.Fatal)
	}
}

//综合服务（百度）（考虑到今后可能接入百度的相关服务）

// 百度结构体
type baiDu struct {
	BaiduFilter baiduFilter //百度过滤服务
}

// 百度过滤服务结构体
type baiduFilter struct {
	APIKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// 初始化百度配置
func (config *Config) initBaidu() {
	viper.SetConfigName("baidu.toml")       // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/app") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read baidu config fatal:"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Application.BaiDu); e != nil {
		err.Error(errors.New("unmarshal baidu config fatal:"+e.Error()), err.Fatal)
	}
}

//LLM相关

//初始化OpenAI配置

// OpenAI结构体
type openai struct {
	ApiKey           string  `mapstructure:"api_key"`           //API密钥
	Model            string  `mapstructure:"model"`             //模型
	Temperature      float64 `mapstructure:"temperature"`       //温度
	TopP             float64 `mapstructure:"top_p"`             //top_p
	MaxTokens        int     `mapstructure:"max_tokens"`        //最大token
	Stop             string  `mapstructure:"stop"`              //停止标志
	PresencePenalty  float64 `mapstructure:"presence_penalty"`  //presence_penalty
	FrequencyPenalty float64 `mapstructure:"frequency_penalty"` //frequency_penalty
}

func (config *Config) initOpenai() {
	viper.SetConfigName("openai.toml")      // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/app") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read openai config fatal:"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Application.Openai); e != nil {
		err.Error(errors.New("unmarshal openai config fatal:"+e.Error()), err.Fatal)
	}
}

//记忆相关

// piencone结构体
type pinecone struct {
	ApiKey string `mapstructure:"api_key"` //API密钥
	Url    string `mapstructure:"url"`     //url
}

// 初始化piencone配置
func (config *Config) initPinecone() {
	viper.SetConfigName("pinecone.toml")    // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/app") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read pinecone config fatal:"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Application.Pinecone); e != nil {
		err.Error(errors.New("unmarshal pinecone config fatal:"+e.Error()), err.Fatal)
	}
}

// 字典名称及字典相关配置
type dict struct {
	FilterDictName string `mapstructure:"filter_dict_name"` //过滤字典名称
	MoodDictName   string `mapstructure:"mood_dict_name"`   //情绪字典名称
}

// 初始化字典配置
func (config *Config) initDict() {
	viper.SetConfigName("dict.toml")        // 配置文件名称(无扩展名)
	viper.SetConfigType("toml")             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config/cfg/app") // 查找配置文件所在的路径
	if e := viper.ReadInConfig(); e != nil {
		err.Error(errors.New("read dict config fatal:"+e.Error()), err.Fatal)
	}
	if e := viper.Unmarshal(&config.Application.Dict); e != nil {
		err.Error(errors.New("unmarshal dict config fatal:"+e.Error()), err.Fatal)
	}
}
