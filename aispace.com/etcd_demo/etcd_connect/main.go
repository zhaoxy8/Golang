package main

import (
        "context"
        "fmt"
        "time"

        "go.etcd.io/etcd/clientv3"
)

func main() {
        cli, err := clientv3.New(clientv3.Config{
                Endpoints:   []string{"161.189.48.115:2379"},
                DialTimeout: 5 * time.Second,
        })
        if err != nil {
                fmt.Printf("connect to etcd failed,err：%v\n", err)
                return
        }
        fmt.Println("connect to etcd success")
        defer cli.Close()
        //put
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
        // value := `[{"path":"/var/log/containers","topic":"dockerlog"},{"path":"/var/log/syslog","topic":"syslog"}]`
        value := `[{"path":"/var/log/containers/*.log","topic":"dockerlog"},{"path":"/var/log/syslog","topic":"syslog"},{"path":"/var/log/messages","topic":"syslog"}]`
        //value := `[{"path":"/var/log/containers/*.log","topic":"dockerlog"},{"path":"/var/log/messages","topic":"syslog"}]`
        _, err = cli.Put(ctx, "/logpath", value)
        cancel()
        if err != nil {
                fmt.Printf("put to etcd failed, err:%v\n", err)
                return
        }
        //get
        ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
        resp, err := cli.Get(ctx, "/logpath")
        if err != nil {
                fmt.Printf("get from etcd failed, err:%v\n", err)
                return
        }
        for _, ev := range resp.Kvs {
                fmt.Printf("%s:%s\n", ev.Key, ev.Value)
        }

}