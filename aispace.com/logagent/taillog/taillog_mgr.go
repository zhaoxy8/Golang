package taillog

import (
	"context"
	"fmt"
	"time"
)

//TailTaskMgr 存储每个path 和 对应的tailtask读取日志任务
type TailTaskMgr struct {
	Path     string
	Topic    string
	TailTask *TailTask
}

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
			fmt.Printf("%s line:%s\n", ttm.Path, line.Text)
		// kafka.SendMsg(cfg.KafkaConf.Topic, line.Text)
		case <-ctx.Done(): // 等待上级通知
			break
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
