package taillog

import (
	"fmt"
	"mymod/logagent/etcd"
	"time"
)

var takMgr *tailLogMgr

// tailTask管理者
type tailLogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry, taskMax int) {
	takMgr = &tailLogMgr{
		logEntry:    logEntryConf, // 把当前日志收集项配置信息保存起来
		tskMap:      make(map[string]*TailTask, taskMax),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲区的通道
	}
	for _, c := range logEntryConf {
		// c: *etcd.LogEntry
		tailObj := NewTailTask(c.Path, c.Topic)
		key := fmt.Sprintf("%s_%s", c.Path, c.Topic)
		takMgr.tskMap[key] = tailObj
	}

	go takMgr.run()
}

// 监听通道有无新配置加入，然后做对应处理
// 配置新增，删除，修改
func (t *tailLogMgr) run() {
	for {
		select {
		case newConfs := <-t.newConfChan:
			// 新增
			newTaskMap := make(map[string]bool) //用一个map把新增项存起来，用来判断是否要删除旧项
			for _, newConf := range newConfs {
				key := fmt.Sprintf("%s_%s", newConf.Path, newConf.Topic)
				_, ok := takMgr.tskMap[key]
				if !ok {
					// 如果原来没有就需要新增
					newTailTask := NewTailTask(newConf.Path, newConf.Topic)
					takMgr.tskMap[key] = newTailTask
				}
				newTaskMap[key] = true
			}
			for _, v := range takMgr.logEntry {
				key := fmt.Sprintf("%s_%s", v.Path, v.Topic)
				if !newTaskMap[key] {
					fmt.Println("删除", key)
					takMgr.tskMap[key].cancelFunc()
					delete(takMgr.tskMap, key)
				}
			}

			fmt.Println("新配置来了！", newConfs)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 向外暴露tskMgr的newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return takMgr.newConfChan
}
