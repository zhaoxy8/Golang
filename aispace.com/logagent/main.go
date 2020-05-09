package main

import (
	"fmt"
	"time"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
)

var (
	cfg = new(conf.Conf)
)

func run() {
	for true {
		select {
		case line := <-taillog.ReadChan():
			fmt.Println("line:", line.Text)
			kafka.SendMsg(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	err := init.MapTo(cfg, "./conf.ini")
	if err != nil {
		fmt.Printf("load conf.init failed", err)
	}
	taillog.Init(cfg.TaillogConf.Path)
	fmt.Println("taillog init success")
	kafka.Init(cfg.KafkaConf.Address)
	fmt.Println("kafka init success")
	run()
}
