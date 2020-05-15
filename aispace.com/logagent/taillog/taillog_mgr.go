package taillog

import (
	"Golong/aispace.com/logagent/etcd"
)

var tailtaskmgr *TailTaskMgr

//TailTaskMgr 存储每个path 和 对应的tailtask读取日志任务
type TailTaskMgr struct {
	LogEntry []*etcd.LogEntry
	TailTask map[string]*TailTask
}

//Init 初始化TailTaskMgr 并执行taillog进行日志操作
func Init(logEntry []*etcd.LogEntry) {
	tailtaskmgr = &TailTaskMgr{
		LogEntry: logEntry,
	}
	taillogmap := make(map[string]*TailTask, len(logEntry))
	for _, v := range logEntry {
		// fmt.Println(v.Path, v.Topic)
		// ctx, cancel := context.WithCancel(context.Background())
		tailtask := NewTailTask(v.Path, v.Topic)
		taillogmap[v.Path] = tailtask
		// defer cancel()
	}
	tailtaskmgr.TailTask = taillogmap
}
