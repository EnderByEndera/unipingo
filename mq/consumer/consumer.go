package consumer

import (
	"log"

	"github.com/streadway/amqp"
)

func Main2() {
	//初始化一个Rabbimtq连接，后跟ip，user，password
	conn, err := amqp.Dial("amqp://melodie:melodie-test@localhost:5672/")
	if err != nil {
		return
	}
	defer conn.Close()
	//创建一个channel的套接字连接
	ch, _ := conn.Channel()

	msgs, err := ch.Consume(
		"python-test", // 队列名
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // 不等待
		nil,           // args
	)
	//定义一个forever，让他驻留在后台，等待消息，来了就消费
	forever := make(chan bool)

	//执行一个go func 完成任务消费
	go func() {
		for d := range msgs {
			//打印body
			log.Printf("received message: '%s'", d.Body)
		}
	}()
	<-forever
}
