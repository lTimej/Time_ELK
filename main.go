package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
}

type Config struct {
	KafkaConfig `ini:"kafka"`
}

func main() {
	config := new(Config)
	err := ini.MapTo(config, "./config/config.ini")
	if err != nil {
		logrus.Println("加载配置失败")
		return
	}

}
