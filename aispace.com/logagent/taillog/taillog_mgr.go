package taillog

import (
	"context"
	"time"
)

//TailTaskMgr 存储每个path 和 对应的tailtask读取日志任务
type TailTaskMgr struct {
	Path     string
	Topic    string
	TailTask *TailTask
}

//LogTopic 日志结构体 存储log和topic
type LogTopic struct {
	Topic string
	Line  string
}

//LogChan channel 用于缓冲日志数据 做一个日志结构体
var LogChan = make(chan *LogTopic, 1000)

//NewTailTaskMgr 构造方法
func NewTailTaskMgr(ctx context.Context, tailtask *TailTask) *TailTaskMgr {
	tailtaskmgr := &TailTaskMgr{
		Path:     tailtask.Path,
		Topic:    tailtask.Topic,
		TailTask: tailtask,
	}
	go tailtaskmgr.Run(ctx)
	return tailtaskmgr
}

//Run 从文件对象中读取数据返回只读chan Line //Intance = tailobj
func (ttm *TailTaskMgr) Run(ctx context.Context) {
	for true {
		select {
		case line := <-ttm.TailTask.Intance.Lines:
			// fmt.Printf("%s line:%s\n", ttm.Path, line.Text)
			logtopic := &LogTopic{
				Topic: ttm.Topic,
				Line:  line.Text,
			}
			//构造一个channel 用于缓冲日志数据 做一个日志结构体
			LogChan <- logtopic
		// kafka.SendMsg(cfg.KafkaConf.Topic, line.Text)
		case <-ctx.Done(): // 等待上级通知
			break
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
