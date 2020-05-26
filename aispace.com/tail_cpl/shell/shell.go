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
		case s3 := <-taillog.GetS3Chan():
			fmt.Printf("需要处理的capsule:%s\n", *s3)
			go run(*s3)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func run(s3 string) {
	//'s3://bixby-submissions/prd/live/capsules/master/bixby.rideShareResolver_KR/0.3.0/capsule.tgz'
	capsule := strings.Split(s3, "/")
	command := "submissions-man-sync.sh" + capsule[7] + "/" + capsule[8]
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
		return
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	S3Path := taillog.GetS3Path()
	delete(S3Path, s3)
}
