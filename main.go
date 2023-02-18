package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/etcd"
	"liujun/Time_ELK/kafka"
	"liujun/Time_ELK/tailf"
	"strings"
)

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int    `ini:"chan_size"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

type TailConfig struct {
	LogfilePath string `ini:"logfile_path"`
}

type Config struct {
	KafkaConfig `ini:"kafka"`
	TailConfig  `ini:"tail"`
	EtcdConfig  `ini:"etcd"`
}

func main() {
	config := new(Config)
	err := ini.MapTo(config, "./config/config.ini")
	if err != nil {
		logrus.Println("加载配置失败")
		return
	}
	err = kafka.Init(strings.Split(config.KafkaConfig.Address, ","), config.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Println("kafka初始化失败,err:", err)
		return
	}
	logrus.Println("kafka初始化成功")
	err = etcd.Init([]string{config.EtcdConfig.Address})
	if err != nil {
		logrus.Println("etcd初始化失败,err:", err)
		return
	}
	logrus.Println("etcd初始化成功")

	etcd_value, err := etcd.GetEtcdValue(config.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Println("获取etcd值失败,err:", err)
		return
	}
	logrus.Println(etcd_value)
	err = tailf.Init(etcd_value)
	if err != nil {
		logrus.Println("tail初始化失败,err:", err)
		return
	}
	logrus.Println("tail初始化成功")
	run()
}
func run() {
	select {}
}
