package config

//配置文件导入引擎

func GetConfig() Config {
	cfg := initConfig()
	return cfg
}
