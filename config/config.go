package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Uid   string `mapstructure:"uid"`  // 用户uid
	Page  int    `mapstructure:"page"` // 动态列表页数
	CType string `mapstructure:"type"`
	Oid   string `mapstructure:"oid"`
}

// 读取配置文件
func Cfg() *Config {
	cfg := &Config{}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Read config error")
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal("Unmarshal config error")
	}
	log.Printf("Read config OK\n")
	return cfg
}
