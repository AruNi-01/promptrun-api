package kafka2

import (
	"github.com/segmentio/kafka-go"
	"os"
)

var KafkaWriter *kafka.Writer

// InitKafkaWriter 初始化 Kafka Writer
func InitKafkaWriter() {
	if KafkaWriter == nil {
		broker := os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_PORT")

		KafkaWriter = &kafka.Writer{
			Addr:                   kafka.TCP(broker),
			Balancer:               &kafka.LeastBytes{}, // 指定分区的 balancer 模式为最小字节分布
			AllowAutoTopicCreation: true,                // 自动创建不存在的 topic
		}
	}
}
