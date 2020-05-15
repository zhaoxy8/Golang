package taillog

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

//TailTask 存储每个tailobj的结构体 tailobj真正打开文件去读取日志
type TailTask struct {
	Path    string
	Topic   string
	Intance *tail.Tail
}

//LogTopic 日志结构体 存储log和topic
type LogTopic struct {
	Topic string
	Line  string
}

//LogChan channel 用于缓冲日志数据 做一个日志结构体
var LogChan = make(chan *LogTopic, 1000)

//NewTailTask 构造方法
func NewTailTask(path string, topic string) *TailTask {
	tailtask := &TailTask{
		Path:  path,
		Topic: topic,
	}
	err := tailtask.init(path)
	if err != nil {
		fmt.Println("NewTailTask err:", err)
	}
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
	go t.Run()
	return
}

//Run 从文件对象中读取数据返回只读chan Line //Intance = tailobj
func (t *TailTask) Run() {
	for true {
		select {
		case line := <-t.Intance.Lines:
			// fmt.Printf("%s line:%s\n", ttm.Path, line.Text)
			logtopic := &LogTopic{
				Topic: t.Topic,
				Line:  line.Text,
			}
			//构造一个channel 用于缓冲日志数据 做一个日志结构体
			LogChan <- logtopic
		// kafka.SendMsg(cfg.KafkaConf.Topic, line.Text) //函数调用函数
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
