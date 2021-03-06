package taillog

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

//var tailtask *TailTask
var logger = log.New(os.Stdout, "[TAIL]", log.Lshortfile|log.Ldate|log.Ltime)

//TailTask 存储每个tailobj的结构体 tailobj真正打开文件去读取日志
type TailTask struct {
	//S3Path 存储s3路径 true
	S3Path map[string]bool
	//S3Chan shell 多协程使用
	S3Chan  chan *string
	tailObj *tail.Tail
}

//NewTailTask 构造方法
func NewTailTask(path string) *TailTask {
	tailTask := &TailTask{
		S3Path: make(map[string]bool, 100),
		S3Chan: make(chan *string),
	}
	err := tailTask.init(path)
	if err != nil {
		logger.Printf("NewTailTask err:%v\n", err)
	}
	return tailTask
}

//Init taillog处理方法
func (t *TailTask) init(path string) (err error) {
	tailObj, err := tail.TailFile(path, tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	})
	if err != nil {
		// fmt.Println("tail file err:", err)
		logger.Println("tail file err:", err)
		return
	}
	t.tailObj = tailObj
	go t.run()
	return
}

func (t *TailTask) run() {
	for true {
		msg, ok := <-t.tailObj.Lines
		if !ok {
			logger.Printf("tail file close reopen, filename: %s\n", t.tailObj.Filename)
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
				//logger.Println("msg:", s3)
			}
		}
	}
}

//GetS3Chan 一个函数，向外暴露只读S3 channel  S3Chan
func  (t *TailTask)GetS3Chan() <-chan *string {
	return t.S3Chan
}

//GetS3Path 一个函数，向外暴露S3路径map  S3Path
func (t *TailTask)GetS3Path() map[string]bool {
	return t.S3Path
}
