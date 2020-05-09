package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var client sarama.SyncProducer

//Init 初始化kafka连接的client
func Init(hosts []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 连接kafka
	client, err = sarama.NewSyncProducer(hosts, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	return
}

//SendMsg 向topic中发送数据
func SendMsg(topic, line string) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(line)
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
