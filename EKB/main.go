package main

import (
	"liujun/Time_ELK/EKB/es"
	"liujun/Time_ELK/EKB/kafka"
	"strings"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int    `ini:"chan_size"`
}
type EsConfig struct {
	Address string `ini:"address"`
	Index   string `ini:"index"`
}

type Config struct {
	KafkaConfig `ini:"kafka"`
	EsConfig    `ini:"es"`
}

var (
	Cfg *Config
)

func main() {
	Cfg = new(Config)
	err := ini.MapTo(Cfg, "./config/config.ini")
	if err != nil {
		logrus.Println("加载配置失败")
		return
	}
	err = kafka.Init(strings.Split(Cfg.KafkaConfig.Address, ","), Cfg.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Println("kafka初始化失败,err:", err)
		return
	}
	logrus.Println("kafka初始化成功")

	go kafka.GetPartition(Cfg.KafkaConfig.Topic)

	err = es.Init(Cfg.EsConfig.Address)
	if err != nil {
		logrus.Println("es初始化失败,err:", err)
		return
	}
	logrus.Println("es初始化成功")
	go es.Put(Cfg.EsConfig.Index, Cfg.KafkaConfig.Topic)
	select {}
}
