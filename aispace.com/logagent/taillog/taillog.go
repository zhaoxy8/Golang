package taillog

import (
	"fmt"

	"github.com/hpcloud/tail"
)

z//TailTask 存储每个tailobj的结构体 tailobj真正打开文件去读取日志
type TailTask struct {
	Path    string
	Topic   string
	Intance *tail.Tail
}

//NewTailTask 构造方法
func NewTailTask(path string, topic string) *TailTask {
	tailtask := &TailTask{
		Path:  path,
		Topic: topic,
	}
	tailtask.init(path)
	return tailtask
}

//init 创建日志文件管理任务对象
func (t *TailTask) init(filename string) (err error) {
	tailObj, err := tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	t.Intance = tailObj
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	return
}

//ReadChan 从文件对象中读取数据返回一个Line结构体型只读chan Line
func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.Intance.Lines
}
