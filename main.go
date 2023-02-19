package main

import (
	"encoding/json"
	"flag"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"liujun/Time_ELK/etcd"
	"liujun/Time_ELK/kafka"
	"liujun/Time_ELK/tail_task_mgr"
	"strings"
)

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int    `ini:"chan_size"`
}

type EtcdConfig struct {
	Address      string `ini:"address"`
	CollectKey   string `ini:"collect_key"`
	CollectValue string `int:"collect_value"`
}

type TailConfig struct {
	LogfilePath string `ini:"logfile_path"`
}

type Config struct {
	KafkaConfig `ini:"kafka"`
	TailConfig  `ini:"tail"`
	EtcdConfig  `ini:"etcd"`
}

var (
	SetEtcd string
	Cfg     *Config
)

func main() {
	flag.StringVar(&SetEtcd, "s", "", "set etcd value by set")
	flag.Parse()
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
	err = etcd.Init([]string{Cfg.EtcdConfig.Address})
	if err != nil {
		logrus.Println("etcd初始化失败,err:", err)
		return
	}
	logrus.Println("etcd初始化成功")
	if SetEtcd != "" {
		data, err := NewEtcdValue()
		if err != nil {
			logrus.Println("设置etcd值序列化失败")
			return
		}
		err = etcd.PutEtcValue(Cfg.EtcdConfig.CollectKey, data)
		if err != nil {
			logrus.Println("设置etcd值失败,err:", err)
			return
		}
	}
	etcd_value, err := etcd.GetEtcdValue(Cfg.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Println("获取etcd值失败,err:", err)
		return
	}
	logrus.Println(etcd_value)

	go etcd.Watch(Cfg.EtcdConfig.CollectKey)

	err = tail_task_mgr.Init(etcd_value)
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

func NewEtcdValue() (string, error) {
	paths := strings.Split(Cfg.TailConfig.LogfilePath, ",")
	var res []map[string]string
	for _, path := range paths {
		ret := map[string]string{"path": path, "topic": Cfg.KafkaConfig.Topic}
		res = append(res, ret)
	}
	data, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
