package main

import (
	"github.com/spf13/viper"
	"log"
)

type config struct {
	CType string `mapstructure:"type"`
	Oid   string `mapstructure:"oid"`
}

var (
	cfg *config
)

func init() {
	cfg = &config{}
}

// 读取配置文件
func getCfg() {
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
}
