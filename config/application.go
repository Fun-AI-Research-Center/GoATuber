package config

import (
	"errors"

	"GoATuber-2.0/err"

	"github.com/spf13/viper"
)

//初始化应用配置

//综合服务（azure）

// Azure azure结构体
type Azure struct {
	AzureOpenai azureOpenai `json:"azure_openai"` //azure openai,文本生成服务
	AzureTTS    azureTTS    `json:"azure_tts"`    //azure tts,文本转语音服务
	AzureSTT    azureSTT    `json:"azure_stt"`    //azure stt,语音转文本服务
}

// azure openai结构体
type azureOpenai struct {
	ApiKey           string  `json:"api_key" mapstructure:"ApiKey"`                     //API密钥
	EndPoint         string  `json:"end_point" mapstructure:"EndPoint"`                 //终结点，在azure openai控制台查询
	DeploymentId     string  `json:"deployment_id" mapstructure:"DeploymentId"`         //模型部署的名称
	ApiVersion       string  `json:"api_version" mapstructure:"ApiVersion"`             //要用于此操作的 API 版本
	Model            string  `json:"model" mapstructure:"Model"`                        //模型
	Temperature      float64 `json:"temperature" mapstructure:"Temperature"`            //温度
	TopP             float64 `json:"top_p" mapstructure:"TopP"`                         //top_p
	MaxTokens        int     `json:"max_tokens" mapstructure:"MaxTokens"`               //最大token
	Stop             string  `json:"stop" mapstructure:"Stop"`                          //停止标志
	PresencePenalty  float64 `json:"presence_penalty" mapstructure:"PresencePenalty"`   //presence_penalty
	FrequencyPenalty float64 `json:"frequency_penalty" mapstructure:"FrequencyPenalty"` //frequency_penalty

	Prompt string `json:"prompt" mapstructure:"Prompt"` //提示(前端将提示句按照特殊规则拼装为字符串，后端再解析)

	MemoryAndClean memoryAndClean `json:"memory_and_clean"` //记忆和清理
	AzureEmbedding azureEmbedding `json:"azure_embedding"`  //azure embedding,嵌入服务
}

// azure embedding azure-openai嵌入服务
type azureEmbedding struct {
	DeploymentID string `json:"deployment_id" mapstructure:"DeploymentID"` //模型部署的名称
	ApiVersion   string `json:"api_version" mapstructure:"ApiVersion"`     //要用于此操作的 API 版本
}

// azure tts结构体
type azureTTS struct {
	ApiKey   string `json:"api_key" mapstructure:"ApiKey"`     //API密钥
	EndPoint string `json:"end_point" mapstructure:"EndPoint"` //终结点，在azure tts控制台查询
	Speak    struct {
		Version    string `json:"version" mapstructure:"Version"`       //指示用于解释文档标记的 SSML 规范的版本。 当前版本为“1.0”。
		XmlLang    string `json:"xmlLang" mapstructure:"XmlLang"`       //根文档的语言。 该值可以包含语言代码，例如 en（英语），也可以包含区域设置，例如 en-US（美国英语）。
		XmlnsMstts string `json:"xmlnsMstts" mapstructure:"XmlnsMstts"` //xmlns:mstts
		Xmlns      string `json:"xmlns" mapstructure:"Xmlns"`           //xmlns
	} `json:"speak"`
	Voice struct {
		Name   string `json:"name" mapstructure:"Name"`     //用于文本转语音输出的语音。 有关受支持的预生成语音的完整列表，请参阅语言支持(https://learn.microsoft.com/zh-cn/azure/cognitive-services/speech-service/language-support?tabs=tts)。
		Effect string `json:"effect" mapstructure:"Effect"` //音频效果处理器，用于在设备上针对特定方案优化合成语音输出的质量。 可选。
		Rate   string `json:"rate" mapstructure:"Rate"`     //语速，填写比默认值高或低的百分比
		Volume string `json:"volume" mapstructure:"Volume"` //音量，填写比默认值高或低的百分比
	} `json:"voice"`
	Authentication string //令牌
	Url            string //请求地址
}

// azure stt结构体
type azureSTT struct {
	ApiKey       string `json:"api_key" mapstructure:"ApiKey"`             //API密钥
	SpeechRegion string `json:"speech_region" mapstructure:"SpeechRegion"` //语音区域
	Language     string `json:"language" mapstructure:"Language"`          //语言
	Format       string `json:"format" mapstructure:"Format"`              //格式
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

// XunFei 讯飞结构体
type XunFei struct {
	XunFeiSpark xunFeiSpark `json:"xun_fei_spark"` //讯飞 星火大模型,文本生成服务
	XunFeiTTS   xunFeiTTS   `json:"xun_fei_tts"`   //讯飞 tts,文本转语音服务
}

type xunFeiSpark struct {
	HostURL   string `json:"host_url" mapstructure:"HostURL"`     //请求地址。如ws(s)://spark-api.xf-yun.com/v2.1/chat
	AppID     string `json:"app_id" mapstructure:"AppID"`         //应用ID
	ApiKey    string `json:"api_key" mapstructure:"ApiKey"`       //API密钥
	ApiSecret string `json:"api_secret" mapstructure:"ApiSecret"` //API密钥

	//以下为请求参数
	Domain      string  `json:"domain" mapstructure:"Domain"`           //领域
	Temperature float64 `json:"temperature" mapstructure:"Temperature"` //温度
	MaxTokens   int     `json:"max_tokens" mapstructure:"MaxTokens"`    //最大token
	TopK        int     `json:"top_k" mapstructure:"TopK"`              //top_k
	Auditing    string  `json:"auditing" mapstructure:"Auditing"`       //审核方式。有没有这个字段我暂且蒙古。

	Prompt string `json:"prompt" mapstructure:"Prompt"` //提示(前端将提示句按照特殊规则拼装为字符串，后端再解析)

	MemoryAndClean memoryAndClean `json:"memory_and_clean"` //记忆和清理
}

type xunFeiTTS struct {
	ApiKey      string `json:"api_key" mapstructure:"ApiKey"`       //API密钥
	ApiSecret   string `json:"api_secret" mapstructure:"ApiSecret"` //API密钥
	AppID       string `json:"app_id" mapstructure:"AppID"`         //API密钥
	XunFeiVoice struct {
		Aue    string `json:"aue" mapstructure:"Aue"`       //音频编码
		Sfl    int    `json:"sfl" mapstructure:"Sfl"`       //流式返回
		Auf    string `json:"auf" mapstructure:"Auf"`       //音频采样率
		Vcn    string `json:"vcn" mapstructure:"Vcn"`       //发音人
		Speed  int    `json:"speed" mapstructure:"Speed"`   //语速
		Volume int    `json:"volume" mapstructure:"Volume"` //音量
		Pitch  int    `json:"pitch" mapstructure:"Pitch"`   //音高
		Reg    string `json:"reg" mapstructure:"Reg"`       //英文发音方式
		Rdn    string `json:"rdn" mapstructure:"Rdn"`       //数字发音方式
	} `json:"xun_fei_voice"` //讯飞语音
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

// BaiDu 百度结构体
type BaiDu struct {
	BaiduFilter baiduFilter `json:"baidu_filter"` //百度过滤服务
	BaiduErnie  baiduErnie  `json:"baidu_ernie"`
}

// 百度 文心一言服务
type baiduErnie struct {
	APPID     string `json:"app_id" mapstructure:"APPID"`
	APIKey    string `json:"api_key" mapstructure:"APIKey"`
	SecretKey string `json:"secret_key" mapstructure:"SecretKey"`

	//Stream bool `json:"stream" mapstructure:"Stream"`//是否使用流式接口（这个疑似是我的事
	Temperature  float64 `json:"temperature" mapstructure:"Temperature"`    //默认0.95。数值越高越随机。范围0到1。
	TopP         float64 `json:"top_p" mapstructure:"TopP"`                 //取值越大，生成的文本多样性越高。范围0到1.0，默认0.8。建议top_p和温度不同时更改。
	PenaltyScore float64 `json:"penalty_score" mapstructure:"PenaltyScore"` //通过对已生成的token增加惩罚，减少重复生成的现象。范围1到2，默认为1。

	Prompt string `json:"prompt" mapstructure:"Prompt"` //提示(前端将提示句按照特殊规则拼装为字符串，后端再解析)

	MemoryAndClean memoryAndClean `json:"memory_and_clean"` //记忆和清理

	AccessToken string `json:"access_token"` //并不是配置文件
	ExpiresIn   int    `json:"expires_in"`   //并不是配置文件
}

// 百度过滤服务结构体
type baiduFilter struct {
	APIKey      string `json:"api_key" mapstructure:"APIKey"`
	SecretKey   string `json:"secret_key" mapstructure:"SecretKey"`
	AccessToken string `json:"access_token"` //并不是配置文件
	ExpiresIn   int    `json:"expires_in"`   //并不是配置文件
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

// Openai OpenAI结构体
type Openai struct {
	ApiKey          string  `json:"api_key" mapstructure:"api_key"`                  //API密钥
	Model           string  `json:"model" mapstructure:"model"`                      //模型
	Temperature     float64 `json:"temperature" mapstructure:"temperature"`          //温度
	TopP            float64 `json:"topP" mapstructure:"top_p"`                       //top_p
	MaxTokens       int     `json:"maxTokens" mapstructure:"max_tokens"`             //最大token
	Stop            string  `json:"stop" mapstructure:"stop"`                        //停止标志
	PresencePenalty float64 `json:"presencePenalty" mapstructure:"presence_penalty"` //presence_penalty

	Prompt string `json:"prompt" mapstructure:"Prompt"` //提示(前端将提示句按照特殊规则拼装为字符串，后端再解析)

	FrequencyPenalty float64        `json:"frequencyPenalty" mapstructure:"frequency_penalty"` //frequency_penalty
	MemoryAndClean   memoryAndClean `json:"memory_and_clean"`                                  //记忆相关
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

// Pinecone piencone结构体
type Pinecone struct {
	ApiKey string `json:"api_key" mapstructure:"api_key"` //API密钥
	Url    string `json:"url" mapstructure:"url"`         //url
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

// Dict 字典名称及字典相关配置
type Dict struct {
	FilterDictName string `json:"filter_dict_name" mapstructure:"filter_dict_name"` //过滤字典名称
	MoodDictName   string `json:"mood_dict_name" mapstructure:"mood_dict_name"`     //情绪字典名称
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

// 一些通配
// 记忆和清理（用于openai和azure openai的gpt服务）
type memoryAndClean struct {
	//短期记忆token上限
	ShortMemoryTokenLimit int `json:"short_memory_token_limit" mapstructure:"ShortMemoryTokenLimit"` //短期记忆token清理限制

	//清理方式
	CleanAll bool `json:"clean_all" mapstructure:"CleanAll"` //清理所有记忆
	CleanOne bool `json:"clean_one" mapstructure:"CleanOne"` //清理一个记忆
}
