package main

import (
	"fmt"
	"time"

	"aispace.com/logagent/kafka"
	"aispace.com/logagent/taillog"
)

func run() {
	for true {
		select {
		case line := <-taillog.ReadChan():
			fmt.Println("line:", line.Text)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
func main() {
	taillog.Init("./my.log")
	fmt.Println("taillog init success")
	kafka.Init([]string{"11.81.1.194:9092", "11.81.1.46:9092"})
	fmt.Println("kafka init success")
	run()
}
