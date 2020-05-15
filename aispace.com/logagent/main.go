package main

import (
	"context"
	"fmt"
	"sync"

	"Golong/aispace.com/logagent/conf"
	"Golong/aispace.com/logagent/etcd"
	"Golong/aispace.com/logagent/kafka"
	"Golong/aispace.com/logagent/taillog"
	"github.com/hpcloud/tail"
	"gopkg.in/ini.v1"
)

//初始化一个Conf结构体指针
var (
	cfg = new(conf.Conf)
	wg  sync.WaitGroup
)

func run(tailpath string, tailline <-chan *tail.Line) {

}

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
	logEntry, err := etcd.GetConf(cfg.EtcdConf.Logkey, cfg.EtcdConf.Timeout)
	fmt.Printf("%#v\n", logEntry)
	// 2.2 派一个哨兵监控日志项的变化，实现热加载

	// 3.使用taillog读取path中的日志发送到kafka
	wg.Add(len(logEntry))
	tailtaskmgrsli := make([]*taillog.TailTaskMgr, len(logEntry))
	for _, v := range logEntry {
		// fmt.Println(v.Path, v.Topic)
		tailtask := taillog.NewTailTask(v.Path, v.Topic)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tailtaskmgr := taillog.NewTailTaskMgr(ctx, tailtask)
		tailtaskmgrsli = append(tailtaskmgrsli, tailtaskmgr)
	}
	for _, v := range tailtaskmgrsli {
		fmt.Println(v.Topic)
	}
	wg.Wait()

	// 3.1 构造一个tailobj的切片每个切片做一个goroute去读取数据
	// run()
}
