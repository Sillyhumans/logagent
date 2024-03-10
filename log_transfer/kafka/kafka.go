package kafka

import (
	"fmt"
	"mymod/log_transfer/es"

	"github.com/IBM/sarama"
)

// kafka消费者实例
func Init(address, topic string) {
	consumer, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				// 消息打包成json格式发送
				data := map[string]interface{}{
					"data": string(msg.Value),
				}
				ld := es.EsData{
					DataLog: data,
					Index:   topic,
				}
				es.SendToChan(&ld)
			}
		}(pc)
	}
}
