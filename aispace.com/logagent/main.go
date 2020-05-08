package main

import (
	"fmt"
	"time"

	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
)

func run() {
	for true {
		select {
		case line := <-taillog.ReadChan():
			fmt.Println("line:", line.Text)
			kafka.SendMsg("web_log", line.Text)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
func main() {
	taillog.Init("./my.log")
	fmt.Println("taillog init success")
	kafka.Init([]string{"11.81.1.194:9092"})
	fmt.Println("kafka init success")
	run()
}
