package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var cli *clientv3.Client

//LogEntry etcd中logpath键值对应的value
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

//Init 初始化etcd cli客户端
func Init(hosts []string, timeout int) (err error) {
	// fmt.Println(hosts)
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   hosts,
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed,err：%v\n", err)
		return
	}
	// fmt.Println("connect to etcd success")
	// defer cli.Close()
	return
}

//GetConf 从etcd中根据key获取配置项
func GetConf(key string, timeout int) (logEntry []*LogEntry, err error) {
	//get
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	// fmt.Printf("%v\n", logEntry)
	// fmt.Println(logEntry == nil)
	for _, ev := range resp.Kvs {
		fmt.Printf("[%s]:%s\n", ev.Key, ev.Value)
		//把解析的LogEntry 组合成一个切片slice
		err = json.Unmarshal(ev.Value, &logEntry)
		if err != nil {
			fmt.Printf("json Unmarshal failed, err:%v\n", err)
			return
		}
	}
	return
}

//WatchConf 监控key的变化
func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	//watch key logpath change
	rch := cli.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			//通知别人
			//1.先判断操作的类型
			var newConf []*LogEntry
			if ev.Type != clientv3.EventTypeDelete {
				//如果是删除操作,手工设置空值
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed,err %v", err)
					continue
				}
			}
			fmt.Printf("get new conf:%v\n", newConf)
			newConfCh <- newConf
		}
	}
}
