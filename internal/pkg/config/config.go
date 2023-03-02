package config

import (
	"github.com/spf13/viper"
	"time"
)

func InitConfig() error {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	return err
}

func GetMysqlConfig() string {
	return viper.GetString("mysql.dsn")
}

func GetFetcherTimeout() time.Duration {
	return time.Duration(viper.GetInt("fetcher.timeout")) * time.Millisecond
}

func GetTaskConfigs() []TaskConfig {
	var configs []TaskConfig
	viper.UnmarshalKey("task", &configs)
	return configs
}

func GetWorkerConfig() ServerConfig {
	var config ServerConfig
	viper.UnmarshalKey("workerconfig", &config)
	return config
}

func GetMasterConfig() ServerConfig {
	var config ServerConfig
	viper.UnmarshalKey("masterconfig", &config)
	return config
}
