package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// kafka consumer
func Init() {
	consumer, err := sarama.NewConsumer([]string{"192.168.70.129:19092", "192.168.70.129:19093", "192.168.70.129:19094"}, nil)
	if err != nil {
		fmt.Printf("kafka消费端启动成功, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("time_elk") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("获取分区失败:err%v\n", err)
		return
	}
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("time_elk", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("分区消费获取失败%d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Println(string(msg.Value))
			}
		}(pc)
	}
	select {}
}
