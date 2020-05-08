package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var tailObj *tail.Tail

func Init(filename string) (err error) {
	tailObj, err := tail.TailFile(filename, tail.Config{
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

func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
}
