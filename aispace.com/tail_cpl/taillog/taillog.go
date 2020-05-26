package taillog

import (
	"fmt"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

//Init taillog处理方法
func Init() {
	filename := "/opt/viv/type-server/logs/webservice-current.log"
	tailFile, err := tail.TailFile(filename, tail.Config{
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

	for true {
		msg, ok := <-tailFile.Lines

		if !ok {
			fmt.Printf("tail file close reopen, filename: %s\n", tailFile.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		idx := strings.Index(msg.Text, "Failed to download CAPSULE_SOURCE")
		if idx >= 1 {
			fmt.Println("msg:", msg.Text)
		}
	}
}
