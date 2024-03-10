package main

import (
	"fmt"
	"mymod/log_transfer/conf"
	"mymod/log_transfer/es"
	"mymod/log_transfer/kafka"

	"gopkg.in/ini.v1"
)

// log transfer
// 将日志从kafka取出发往es

func main() {
	// 加载配置文件
	cfg := &conf.LogTransfer{}
	err := ini.MapTo(cfg, "./conf/conf.ini")
	if err != nil {
		fmt.Println("init config failed, err:", err)
		return
	}
	fmt.Println(cfg)

	// 初始化es
	err = es.Init(cfg.ESConf.Address, 10000)
	if err != nil {
		fmt.Println("init es failed, err:", err)
		return
	}
	fmt.Println("init es success!")

	// 初始化kafka
	kafka.Init(cfg.KafkaConf.Address, cfg.KafkaConf.Topic)

	fmt.Println("init kafka success!")
	select {}
	// 向es发送数据
}
