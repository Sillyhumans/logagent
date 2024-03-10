package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	cil *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(address string, timeout time.Duration) (err error) {
	cil, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{address},
		DialTimeout: timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return err
	}
	fmt.Println("connect to etcd success")
	//defer cil.Close()
	return nil
}

// 从etcd中获取配置项
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// get
	resp, err := cil.Get(ctx, key)
	cancel()

	if err != nil {
		fmt.Println("get keys failed", err)
		return
	}
	for _, val := range resp.Kvs {
		err := json.Unmarshal(val.Value, &logEntryConf)
		if err != nil {
			fmt.Println("unmarshal etcd value failed err:", err)
			return nil, err
		}
	}
	return
}

// 设置哨兵监控ectd中的配置项是否有变动
func WatchConf(key string, newLogChan chan<- []*LogEntry) {
	ch := cil.Watch(context.Background(), key)
	// 从通道中尝试取值
	for wresp := range ch {
		for _, ev := range wresp.Events {
			// 通知tailTask
			fmt.Printf("Type: %v key: %v value: %v\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			// 通知taillog.tskMgr
			// 先判断操作的类型
			var newConf []*LogEntry
			if ev.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed, err: %v\n", err)
					continue
				}
				fmt.Printf("get new conf: %v\n", newConf)
			}
			newLogChan <- newConf
		}
	}
}
