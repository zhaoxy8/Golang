package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var client sarama.SyncProducer

func Init(hosts []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 连接kafka
	client, err := sarama.NewSyncProducer(hosts, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
}
