package main
import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)
func main(){
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"11.81.1.164:2379"},
		DialTimeout: 5*time.Second,
	})
}