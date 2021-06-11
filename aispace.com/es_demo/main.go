package main

import (
	"context"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

// 连接ES 并发送日志数据

var logger = log.New(os.Stdout, "[ES]", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Ltime)

//Person ...
type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {

	client, err := elastic.NewClient(elastic.SetURL("es-cn-nif1w339o000rryq1.elasticsearch.aliyuncs.com:9200"))
	if err != nil {
		panic(err)
	}
	logger.Println("connect to es success")
	p1 := Person{
		Name:    "wanglei",
		Age:     32,
		Married: true,
	}
	//链式操作
	put1, err := client.Index().
		Index("go-user").
		BodyJson(p1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	logger.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
