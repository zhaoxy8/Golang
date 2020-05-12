package main

import (
	"fmt"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/etcd"
	"Golong/aispace.com/logagent/kafka"
	"gopkg.in/ini.v1"
)

//初始化一个Conf结构体指针
var (
	cfg = new(conf.Conf)
)

// func run() {
// 	for true {
// 		select {
// 		case line := <-taillog.ReadChan():
// 			fmt.Println("line:", line.Text)
// 			kafka.SendMsg(cfg.KafkaConf.Topic, line.Text)
// 		default:
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}
// }

func main() {
	//1.把ini文件转换为结构体
	err := ini.MapTo(cfg, "./conf.ini")
	if err != nil {
		fmt.Printf("load conf.init failed%v\n", err)
	}
	// err = taillog.Init(cfg.TaillogConf.Path)
	// if err != nil {
	// 	fmt.Printf("taillog init failed%v\n", err)
	// }
	// fmt.Println("taillog init success")
	//2.连接etcd
	err = etcd.Init(cfg.EtcdConf.Address, cfg.EtcdConf.Timeout)
	if err != nil {
		fmt.Printf("etcd init failed%v\n", err)
	}
	fmt.Println("etcd init success")
	err = kafka.Init(cfg.KafkaConf.Address)
	if err != nil {
		fmt.Printf("kafka init failed%v\n", err)
	}
	fmt.Println("kafka init success")
	// 2.1从etcd中获取日志收集项的配置信息
	logEntry, err := etcd.GetConf("/logpath", cfg.EtcdConf.Timeout)
	fmt.Printf("%#v\n", logEntry)
	for _, v := range logEntry {
		fmt.Println(v.Path, v.Topic)
	}
	// 2.2 派一个哨兵监控日志项的变化，实现热加载
	//执行读取日志发送到kakfa
	// run()
}
