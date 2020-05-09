package main

import (
	"fmt"
	"time"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
)

func run(conf *conf.Conf) {
	for true {
		select {
		case line := <-taillog.ReadChan():
			fmt.Println("line:", line.Text)
			kafka.SendMsg(conf.Topic, line.Text)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
func main() {
	conf := conf.NewConf("conf.ini")
	taillog.Init(conf.Path)
	fmt.Println("taillog init success")
	kafka.Init(conf.Hosts)
	fmt.Println("kafka init success")
	run(conf)
}
