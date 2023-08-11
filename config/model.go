package config

//配置文件结构体模型

// 消息监听
type monitor struct {
	BiliBili bool `mapstructure:"bilibili"`
	Other    bool `mapstructure:"other"`
}

// 过滤器
type filter struct {
	Dict  bool `mapstructure:"dict"`  //字典过滤器
	Baidu bool `mapstructure:"baidu"` //百度过滤器
	Other bool `mapstructure:"other"` //其他过滤器
}

// 语言模型配置
type llm struct {
	Openai      bool `mapstructure:"openai"`
	AzureOpenai bool `mapstructure:"azure_openai"`
	Other       bool `mapstructure:"other"`
}

// 语音模型配置
type voice struct {
	Azure bool `mapstructure:"azure"`
	XFyun bool `mapstructure:"xfyun"`
	Other bool `mapstructure:"other"`
}

// 对话配置
type speech struct {
	Azure bool `mapstructure:"azure"`
	Other bool `mapstructure:"other"`
}

// 情感配置
type mood struct {
	Dict  bool `mapstructure:"dict"`
	Other bool `mapstructure:"other"`
}

// 代理配置
type proxy struct {
	UseProxy bool   `mapstructure:"use_proxy"`
	ProxyURL string `mapstructure:"proxy_url"`
}

// 监听包
type listen struct {
	BiliBili biliBili //BiliBili监听
}

// 应用包
type application struct {
	Openai   openai   //OpenAI
	Azure    azure    //Azure
	XunFei   xunFei   //讯飞
	BaiDu    baiDu    //百度
	Pinecone pinecone //Pinecone
	Dict     dict     //字典配置
}

// 工具包
type tool struct {
	Memory memory //记忆相关
}
