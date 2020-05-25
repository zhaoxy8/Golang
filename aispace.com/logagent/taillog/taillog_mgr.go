package taillog

import (
	"fmt"
	"time"

	"Golong/aispace.com/logagent/etcd"
)

var tailtaskmgr *TailTaskMgr

//TailTaskMgr tailtask 管理者
type TailTaskMgr struct {
	LogEntry    []*etcd.LogEntry
	TailTask    map[*etcd.LogEntry]*TailTask
	newConfChan chan []*etcd.LogEntry
}

//Init 初始化TailTaskMgr 并执行taillog进行日志操作
func Init(logEntry []*etcd.LogEntry) {
	tailtaskmgr = &TailTaskMgr{
		LogEntry:    logEntry,
		newConfChan: make(chan []*etcd.LogEntry), //无缓冲区通道
	}
	taillogmap := make(map[*etcd.LogEntry]*TailTask, len(logEntry))
	for _, v := range logEntry {
		// fmt.Println(v.Path, v.Topic)
		// ctx, cancel := context.WithCancel(context.Background())
		// 真正执行tail日志操作，并把数据写入channel中
		tailtask := NewTailTask(v.Path, v.Topic)
		taillogmap[v] = tailtask
		// defer cancel()
	}
	tailtaskmgr.TailTask = taillogmap
	tailtaskmgr.run()
}

func (t *TailTaskMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			//1.配置新增
			//2.配置删除
			//3.配置变更
			fmt.Println("新增配置来了", newConf)
		default:
			time.Sleep(time.Second)
		}
	}

}

//NewConfChan 一个函数，向外暴露TailTaskMgr  newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tailtaskmgr.newConfChan
}
