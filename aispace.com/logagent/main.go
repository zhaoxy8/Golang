package main

import (
	"fmt"
	"sync"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/etcd"
	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
	"gopkg.in/ini.v1"
)

//初始化一个Conf结构体指针
var (
	cfg = new(conf.Conf)
	wg  sync.WaitGroup
)

func main() {
	//1.把ini文件转换为结构体
	err := ini.MapTo(cfg, "./conf.ini")
	if err != nil {
		fmt.Printf("load conf.init failed%v\n", err)
	}
	//2.连接etcd
	err = etcd.Init(cfg.EtcdConf.Address, cfg.EtcdConf.Timeout)
	if err != nil {
		fmt.Printf("etcd init failed%v\n", err)
	}
	fmt.Println("etcd init success")
	err = kafka.Init(cfg.KafkaConf.Address, cfg.KafkaConf.LogSize)
	if err != nil {
		fmt.Printf("kafka init failed%v\n", err)
	}
	fmt.Println("kafka init success")
	// 2.1从etcd中获取日志收集项的配置信息
	logEntry, err := etcd.GetConf(cfg.EtcdConf.Logkey, cfg.EtcdConf.Timeout)
	// fmt.Printf("%#v\n", logEntry)
	// 3.使用taillog读取path中的日志发送到kafka
	taillog.Init(logEntry)
	newConfChan := taillog.NewConfChan()
	wg.Add(1)
	//从taillog包中获取对外暴露的通道
	// 2.2 派一个哨兵监控日志项的变化，实现热加载
	go etcd.WatchConf(cfg.EtcdConf.Logkey, newConfChan) //哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()

}
