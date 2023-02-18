package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/kafka"
	"strings"
)

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
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
	err = kafka.Init(strings.Split(config.KafkaConfig.Address, ","))
	if err != nil {
		logrus.Println("kafka初始化失败,err:", err)
		return
	}
	logrus.Println("kafka初始化成功")
}
