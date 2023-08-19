package config

//配置文件结构体模型

// Monitor 消息监听
type Monitor struct {
	BiliBili bool `json:"bilibili" mapstructure:"bilibili"`
	None     bool `json:"none" mapstructure:"none"`
}

// Filter 过滤器
type Filter struct {
	Dict  bool `json:"dict" mapstructure:"dict"`   //字典过滤器
	Baidu bool `json:"baidu" mapstructure:"baidu"` //百度过滤器
	Other bool `json:"other" mapstructure:"other"` //其他过滤器
}

// LLM 语言模型配置
type LLM struct {
	Openai      bool `json:"openai" mapstructure:"openai"`
	AzureOpenai bool `json:"azure_openai" mapstructure:"azure_openai"`
	Other       bool `json:"other" mapstructure:"other"`
}

// Voice 语音模型配置
type Voice struct {
	Azure bool `json:"azure" mapstructure:"azure"`
	XFyun bool `json:"xfyun" mapstructure:"xfyun"`
	Other bool `json:"other" mapstructure:"other"`
}

// Speech 对话配置
type Speech struct {
	Azure bool `json:"azure" mapstructure:"azure"`
	Other bool `json:"other" mapstructure:"other"`
}

// Mood 情感配置
type Mood struct {
	Dict  bool `json:"dict" mapstructure:"dict"`
	Other bool `json:"other" mapstructure:"other"`
}

// Proxy 代理配置
type Proxy struct {
	UseProxy bool   `json:"use_proxy" mapstructure:"use_proxy"`
	ProxyURL string `json:"proxy_url" mapstructure:"proxy_url"`
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
