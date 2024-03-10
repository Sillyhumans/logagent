package taillog

import (
	"context"
	"fmt"
	"mymod/logagent/kafka"
	"time"

	"github.com/nxadm/tail"
)

// TailTask: 一个日志收集任务
type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	err := tailObj.init()
	if err != nil {
		fmt.Println("tailObj init failed, err:", err)
		return
	}
	return tailObj
}

func (t *TailTask) init() (err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	// 实例化一个tail对象
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("new tail failed, err:", err)
		return
	}
	go t.run()
	return
}

func (t *TailTask) run() {
	for {
		// 获取日志
		select {
		case <-t.ctx.Done():
			fmt.Println("tail task结束了...", t.path, t.topic)
			return
		case line := <-t.instance.Lines:
			// 现将数据存入通道中，再开启一个后台存入kafka

			kafka.SendToChan(t.topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}
