package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// 专门往kafka写日志的模块

type logData struct {
	topic string
	data  string
}

var (
	clinet      sarama.SyncProducer //声明一个全局连接kafka生产者
	logDataChan chan *logData
)

func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

func Init(address []string, chanMax int) (err error) {

	logDataChan = make(chan *logData, chanMax)

	config := sarama.NewConfig()

	// tailf包的使用
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要等待所有响应，即leader和follower都确认

	config.Producer.Partitioner = sarama.NewRandomPartitioner // 重新选出一个 partition

	config.Producer.Return.Successes = true // 成功交付的消息即将在success channel返回

	// 连接kafka
	clinet, err = sarama.NewSyncProducer(address, config)

	if err != nil {
		fmt.Println("producer closed err:", err)
		return err
	}

	// 从通道中取数据发往kafka
	go sendMsg()

	return nil
}

// 真正往kafka发送日志的函数
func sendMsg() (err error) {
	for {
		select {
		case ld := <-logDataChan:
			// 构造一个消息
			msger := &sarama.ProducerMessage{}
			msger.Topic = ld.topic
			fmt.Println(msger.Topic)
			msger.Value = sarama.StringEncoder(ld.data)
			pid, offset, err := clinet.SendMessage(msger)
			if err != nil {
				fmt.Println("send message failed, err:", err)
				return err
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
