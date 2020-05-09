package main

import (
	"fmt"
	"time"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
	"gopkg.in/ini.v1"
)

//初始化一个Conf结构体指针
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
	//1.把ini文件转换为结构体
	err := ini.MapTo(cfg, "./conf.ini")
	if err != nil {
		fmt.Printf("load conf.init failed%v\n", err)
	}
	err = taillog.Init(cfg.TaillogConf.Path)
	if err != nil {
		fmt.Printf("taillog init failed%v\n", err)
	}
	fmt.Println("taillog init success")
	err = kafka.Init(cfg.KafkaConf.Address)
	if err != nil {
		fmt.Printf("kafka init failed%v\n", err)
	}
	fmt.Println("kafka init success")
	//执行读取日志发送到kakfa
	run()
}
