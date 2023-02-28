package config

import (
	"crawler/internal/pkg/spider"
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

func GetTaskConfigs() []spider.TaskConfig {
	var configs []spider.TaskConfig
	viper.UnmarshalKey("task", &configs)
	return configs
}
