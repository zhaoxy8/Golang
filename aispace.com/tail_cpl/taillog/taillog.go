package taillog

import (
	"fmt"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

// var (
// 	S3Path map[string]bool

// 	S3Chan chan string
// )
var tailtask *TailTask

//TailTask 存储每个tailobj的结构体 tailobj真正打开文件去读取日志
type TailTask struct {
	//S3Path 存储s3路径
	S3Path map[string]bool
	//S3Chan shell 多协程使用
	S3Chan  chan *string
	tailobj *tail.Tail
}

//NewTailTask 构造方法
func NewTailTask(path string) *TailTask {
	tailtask = &TailTask{
		S3Path: make(map[string]bool, 100),
		S3Chan: make(chan *string),
	}
	err := tailtask.init(path)
	if err != nil {
		fmt.Println("NewTailTask err:", err)
	}
	return tailtask
}

//Init taillog处理方法
func (t *TailTask) init(path string) (err error) {
	tailobj, err := tail.TailFile(path, tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	t.tailobj = tailobj
	go t.run()
	return
}

func (t *TailTask) run() {
	for true {
		msg, ok := <-t.tailobj.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename: %s\n", t.tailobj.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		idx := strings.Index(msg.Text, "Failed to download CAPSULE_SOURCE")
		if idx >= 1 {
			// 字符串截取s3路径
			s3 := strings.Split(msg.Text, "'")[1]
			// 如果s3路径不在map中就添加到map列表中,同步完成需要删除map中的key
			_, ok := t.S3Path[s3]
			if !ok {
				t.S3Path[s3] = true
				//把s3路径地址放到chan中
				t.S3Chan <- &s3
			}
			fmt.Println("msg:", s3)

		}
	}
}

//NewConfChan 一个函数，向外暴露只读tailtask  S3Chan
func NewConfChan() <-chan *string {
	return tailtask.S3Chan
}
