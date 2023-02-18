package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var (
	Client sarama.SyncProducer
	err    error
)

func Init(addr []string) error {
	fmt.Println(addr)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	Client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		return err
	}
	return nil
}
