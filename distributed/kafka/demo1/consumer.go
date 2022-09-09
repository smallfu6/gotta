package main

/*
	kafka consumer

	创建topic
	kafka-topics.sh --create --topic topic-demo --partitions 3 --bootstrap-server localhost:9092
*/

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}

	// 根据topic取到所有的分区
	partitionList, err := consumer.Partitions("topic-d1")
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)

	var wg sync.WaitGroup

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("topic-d1", int32(partition),
			sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n",
				partition, err)
			return
		}
		defer pc.AsyncClose()

		wg.Add(1)
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n",
					msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}

	wg.Wait()
}
