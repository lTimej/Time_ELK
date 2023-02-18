package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/kafka"
	"liujun/Time_ELK/tailf"
	"strings"
)

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int    `ini:"chan_size"`
}

type TailConfig struct {
	LogfilePath string `ini:"logfile_path"`
}

type Config struct {
	KafkaConfig `ini:"kafka"`
	TailConfig  `ini:"tail"`
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

	err = tailf.Init(config.TailConfig.LogfilePath)
	if err != nil {
		logrus.Println("tail初始化失败,err:", err)
		return
	}
	logrus.Println("tail初始化成功")
}
