package main

import (
	"fmt"
	"time"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
	"gopkg.in/ini.v1"
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
	err := ini.MapTo(cfg, "./conf.ini")
	if err != nil {
		fmt.Printf("load conf.init failed", err)
	}
	err = taillog.Init(cfg.TaillogConf.Path)
	if err != nil {
		fmt.Printf("taillog init failed", err)
	}
	fmt.Println("taillog init success")
	err = kafka.Init(cfg.KafkaConf.Address)
	if err != nil {
		fmt.Printf("kafka init failed", err)
	}
	fmt.Println("kafka init success")
	run()
}
