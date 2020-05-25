package taillog

import (
	"Golong/aispace.com/logagent/etcd"
)

var tailtaskmgr *TailTaskMgr

//TailTaskMgr 存储每个LogEntry 和 对应的tailtask
type TailTaskMgr struct {
	LogEntry []*etcd.LogEntry
	TailTask map[*etcd.LogEntry]*TailTask
}

//Init 初始化TailTaskMgr 并执行taillog进行日志操作
func Init(logEntry []*etcd.LogEntry) {
	tailtaskmgr = &TailTaskMgr{
		LogEntry: logEntry,
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
}
