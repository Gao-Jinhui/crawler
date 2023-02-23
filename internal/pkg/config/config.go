package config

import (
	"github.com/spf13/viper"
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
