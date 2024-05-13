package kafka2

import (
	"context"
	"github.com/segmentio/kafka-go"
	"os"
	"promptrun-api/utils"
)

// SendMessage 同步发送消息
func SendMessage(topic string, key, message string) error {
	if err := KafkaWriter.WriteMessages(context.Background(), kafka.Message{
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
		if err := KafkaWriter.WriteMessages(context.Background(), kafka.Message{
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
		reader := getReaderWithTopic(topic)

		// 读取消息，每次最大 10M
		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				utils.Log().Error("", "read message fail, errMsg: %s", err.Error())
				break
			}
			utils.Log().Info("", "【Consumer】from topic: %s, received message: %s", topic, string(message.Value))
			handler(string(message.Value))
		}
	}(topic, handler)
}

// getReaderWithTopic 获取指定 topic 的 Reader
func getReaderWithTopic(topic string) *kafka.Reader {
	broker := os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_PORT")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	// 设置 offset 为最新，即只消费新消息
	if reader.SetOffset(kafka.LastOffset) != nil {
		utils.Log().Error("", "set offset fail")
	}

	return reader
}
