package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	Consumer sarama.Consumer
	err      error
	KfKMsg   chan string
)

// kafka consumer
func Init(addr []string, chan_size int) error {
	Consumer, err = sarama.NewConsumer(addr, nil)
	if err != nil {
		return err
	}
	KfKMsg = make(chan string, chan_size)
	return nil

}

func GetPartition(topic string) {
	partitionList, err := Consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("获取分区失败:err%v\n", err)
		return
	}
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := Consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("分区消费获取失败%d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		fmt.Println(topic)
		go func(sarama.PartitionConsumer) {
			for {
				for msg := range pc.Messages() {
					fmt.Println(string(msg.Value))
					KfKMsg <- string(msg.Value)
				}
			}
		}(pc)
	}
	select {}
	fmt.Println(11111)
}
