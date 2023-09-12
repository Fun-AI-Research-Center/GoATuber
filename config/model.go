package config

//配置文件结构体模型

// Monitor 消息监听
type Monitor struct {
	BiliBili bool `json:"bilibili" mapstructure:"BiliBili"`
	None     bool `json:"none" mapstructure:"None"`
}

// Filter 过滤器
type Filter struct {
	Dict  bool `json:"dict" mapstructure:"Dict"`   //字典过滤器
	Baidu bool `json:"baidu" mapstructure:"Baidu"` //百度过滤器
	None  bool `json:"none" mapstructure:"None"`   //其他过滤器
}

// LLM 语言模型配置
type LLM struct {
	Openai      bool `json:"openai" mapstructure:"Openai"`
	AzureOpenai bool `json:"azure_openai" mapstructure:"AzureOpenai"`
	XunFeiSpark bool `json:"xun_fei" mapstructure:"XunFeiSpark"`    //讯飞星火大模型
	BaiduErnie  bool `json:"baidu_ernie" mapstructure:"BaiduErnie"` //百度文心一言
	None        bool `json:"none" mapstructure:"None"`
}

// Voice 语音模型配置
type Voice struct {
	Azure bool `json:"azure" mapstructure:"Azure"`
	XFyun bool `json:"xfyun" mapstructure:"XFyun"`
	None  bool `json:"none" mapstructure:"None"`
}

// Speech 对话配置
type Speech struct {
	Azure bool `json:"azure" mapstructure:"Azure"`
	None  bool `json:"none" mapstructure:"None"`
}

// Mood 情感配置
type Mood struct {
	Dict bool `json:"dict" mapstructure:"Dict"`
	None bool `json:"none" mapstructure:"None"`
}

// Proxy 代理配置
type Proxy struct {
	UseProxy bool   `json:"use_proxy" mapstructure:"UseProxy"`
	ProxyURL string `json:"proxy_url" mapstructure:"ProxyURL"`
}

// Listen 监听包
type Listen struct {
	BiliBili BiliBili `json:"bilibili"` //BiliBili监听
}

// 应用包
type application struct {
	Openai   Openai   `json:"openai"`   //OpenAI
	Azure    Azure    `json:"azure"`    //Azure
	XunFei   XunFei   `json:"xunfei"`   //讯飞
	BaiDu    BaiDu    `json:"baidu"`    //百度
	Pinecone Pinecone `json:"pinecone"` //Pinecone
	Dict     Dict     `json:"dict"`     //字典配置
}

// 工具包
type tool struct {
	Memory Memory `json:"memory"` //记忆相关
}
