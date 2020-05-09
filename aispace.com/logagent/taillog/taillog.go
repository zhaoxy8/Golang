package taillog

import (
	"fmt"

	"github.com/hpcloud/tail"
)

var tailObj *tail.Tail

//Init 创建日志文件对象
func Init(filename string) (err error) {
	tailObj, err = tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	return
}

//ReadChan 从文件对象中读取数据返回一个Line结构体型只读channel
func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
}
