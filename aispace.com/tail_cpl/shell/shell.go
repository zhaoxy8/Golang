package shell

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"Golong/aispace.com/tail_cpl/taillog"
)

func Init() {
	for {
		select {
		//获取chan中需要处理的s3,处理完成后堵塞
		case s3 := <-taillog.GetS3Chan():
			fmt.Printf("需要处理的:%s\n", *s3)
			go run(*s3)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func run(s3 string) {
	//'s3://bixby-submissions/prd/live/capsules/master/bixby.rideShareResolver_KR/0.3.0/capsule.tgz'
	capsule := strings.Split(s3, "/")
	command := "submissions-man-sync.sh " + capsule[7] + "/" + capsule[8]
	fmt.Printf("需要执行的命令:%s\n", command)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	S3Path := taillog.GetS3Path()
	fmt.Printf("S3Path map :%v\n", S3Path)
	//处理完成后把map中的key删除掉
	defer delete(S3Path, s3)
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s\n", command, err.Error())
		return
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s\n", command, string(output))
}
