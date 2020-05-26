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
	tailTaskMap map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

//Init 初始化TailTaskMgr 并执行taillog进行日志操作
func Init(logEntry []*etcd.LogEntry) {
	tailtaskmgr = &TailTaskMgr{
		LogEntry:    logEntry,
		tailTaskMap: make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), //无缓冲区通道
	}
	// taillogmap := make(map[*etcd.LogEntry]*TailTask, len(logEntry))
	for _, v := range logEntry {
		// fmt.Println(v.Path, v.Topic)
		// ctx, cancel := context.WithCancel(context.Background())
		// 真正执行tail日志操作，并把数据写入channel中
		mk := fmt.Sprintf("%s_%s", v.Path, v.Topic)
		tailobj := NewTailTask(v.Path, v.Topic)
		tailtaskmgr.tailTaskMap[mk] = tailobj
		// defer cancel()
	}
	// tailtaskmgr.TailTask = taillogmap
	go tailtaskmgr.run()
}

func (t *TailTaskMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("新增配置来了", newConf)
			for _, conf := range newConf {
				//1.配置新增
				//2.配置删除
				//3.配置变更
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				// 如果新配置在老的配置切片中，就不做处理
				_, ok := t.tailTaskMap[mk]
				if ok {
					continue
				}
				// 如果新配置不在老的配置切片中，就新启动tailobj做处理。tailTaskMap中添加任务
				tailobj := NewTailTask(conf.Path, conf.Topic)
				t.tailTaskMap[mk] = tailobj
			}
		default:
			time.Sleep(time.Second)
		}
	}

}

//NewConfChan 一个函数，向外暴露TailTaskMgr  newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tailtaskmgr.newConfChan
}
