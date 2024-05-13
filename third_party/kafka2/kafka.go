package kafka2

import (
	"context"
	"github.com/segmentio/kafka-go"
	"os"
	"promptrun-api/utils"
)

// SendMessage 同步发送消息
func SendMessage(topic string, key, message string) error {
	if _, err := KafkaConn.WriteMessages(kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(message),
	}); err != nil {
		utils.Log().Error("", "send message fail, errMsg: %s", err.Error())
		return err
	}

	return nil
}

// SendMessageAsync 异步发送消息
func SendMessageAsync(topic string, key, message string, onSuccess func(message string), onFail func(message string)) {
	// 创建一个用于通知消息发送结果的 channel
	sendResult := make(chan error, 1)

	// 启动一个 goroutine 异步发送消息
	go func(topic string, key, message string, sendResult chan<- error) {
		if _, err := KafkaConn.WriteMessages(kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: []byte(message),
		}); err != nil {
			utils.Log().Error("", "send message fail, errMsg: %s", err.Error())
			sendResult <- err
		} else {
			sendResult <- nil
		}
	}(topic, key, message, sendResult)

	// 在另一个 goroutine 中等待消息发送结果，并调用相应的回调函数
	go func(sendResult <-chan error) {
		err := <-sendResult
		if err != nil {
			onFail(message)
		} else {
			onSuccess(message)
		}
	}(sendResult)
}

// Subscribe 订阅消息
func Subscribe(topic string, handler func(msg string)) {
	go func(topic string, handler func(msg string)) {
		conn := getConnWithTopic(topic)
		defer func(conn *kafka.Conn) {
			err := conn.Close()
			if err != nil {
				utils.Log().Error("", "close kafka conn fail, errMsg: %s", err.Error())
			}
		}(conn)

		// 读取消息，每次最大 10M
		for {
			message, err := conn.ReadMessage(10e6)
			if err != nil {
				utils.Log().Error("", "read message fail, errMsg: %s", err.Error())
				break
			}
			handler(string(message.Value))
		}
	}(topic, handler)
}

// getConnWithTopic 获取指定 topic 的连接
func getConnWithTopic(topic string) *kafka.Conn {
	broker := os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_PORT")

	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		utils.Log().Panic("", "Kafka Broker 连接失败, errMsg: %s", err.Error())
		panic(err)
	}
	return conn
}
