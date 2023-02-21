package main

import (
	"fmt"
	"time"

	"codefridge/mq/consumer"

	"github.com/streadway/amqp"
)

func main() {
	//初始化一个Rabbimtq连接，后跟ip，user，password
	conn, err := amqp.Dial("amqp://melodie:melodie-test@localhost:5672/")
	if err != nil {
		return
	}
	defer conn.Close()
	//创建一个channel的套接字连接
	ch, _ := conn.Channel()
	//创建一个指定的队列
	_, err = ch.QueueDeclare(
		"python-test", // 队列名
		false,         // durable
		false,         // 不使用删除？
		false,         // exclusive
		false,         // 不必等待
		nil,           // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	//定义上传的消息
	body := "work message"
	//调用Publish上传消息1到指定的work队列当中
	err = ch.Publish(
		"",     // exchange
		"work", // 队列名
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			//[]byte化body
			Body: []byte(body),
		})
	go func() {
		consumer.Main2()
	}()

	time.Sleep(time.Second * 10)
}
