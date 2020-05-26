package shell

import (
	"fmt"
	"time"

	"Golong/aispace.com/tail_cpl/taillog"
)

func Init() {

	for {
		select {
		case s3 := <-taillog.GetS3Chan():
			fmt.Printf("需要处理的capsule:%s\n", *s3)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
