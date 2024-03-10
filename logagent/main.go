package main

import (
	"fmt"
	"mymod/logagent/conf"
	"mymod/logagent/etcd"
	"mymod/logagent/kafka"
	"mymod/logagent/taillog"
	"mymod/logagent/utils"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

// logagent入口

var (
	cfg = new(conf.Appconf)
)

func main() {

	// 加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")

	fmt.Println(cfg)

	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
	}

	// 1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.ChanMax)

	if err != nil {
		fmt.Printf("init kafka failed, err: %v \n", err)
		return
	}

	fmt.Println("init kafka success")

	// 2.初始化ETCD
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)

	if err != nil {
		fmt.Println("init etcd failed, err:", err)
		return
	}

	fmt.Println("init etcd success")

	// 为了实现每个logagent拉取自己独有配置，要义ip作为区分
	// 本地占时用127.0.0.1
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	fmt.Println(ipStr)
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, "127.0.0.1:2379")
	// 从etcd中获取日志收集项的配置信息
	logagent, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Println("etcd.Getconf falied, err", err)
		return
	}
	fmt.Println("get etcd conf success!")

	// 3.初始化日志收集模块 收集日志发送到kafka
	// 3.1 循环每个日志收集项创建对应的tailOBj
	// 3.2 发往kafka
	taillog.Init(logagent, cfg.EtcdConf.TaskMax)

	var wg sync.WaitGroup
	wg.Add(1)
	//派一个哨兵去监控etcd中日志收集项的配置是否有变动
	newConfChan := taillog.NewConfChan() // 从taillog包中获取对外暴露的通道
	go etcd.WatchConf(etcdConfKey, newConfChan)
	wg.Wait()
	//run()
}
