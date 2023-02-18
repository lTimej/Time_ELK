package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client sarama.SyncProducer
	Msg    chan *sarama.ProducerMessage
	err    error
)

func Init(addr []string, capacity int) error {
	fmt.Println(addr)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	Client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		return err
	}
	Msg = make(chan *sarama.ProducerMessage, capacity)
	go SendMsg()
	return nil
}
func SendMsg() {
	for {
		select {
		case msg := <-Msg:
			logrus.Println("=========", msg.Topic, msg.Value, "===========")
			_, _, err := Client.SendMessage(msg)
			if err != nil {
				fmt.Println("发送消息失败s, err:", err)
				return
			}
		}
	}

}
