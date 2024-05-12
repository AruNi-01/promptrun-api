package kafka2

import (
	"context"
	"github.com/segmentio/kafka-go"
	"os"
	"promptrun-api/utils"
	"time"
)

var KafkaConn *kafka.Conn

func InitKafkaBroker() {
	broker := os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_PORT")

	conn, err := kafka.DialContext(context.Background(), "tcp", broker)
	if err != nil {
		utils.Log().Panic("", "Kafka Broker 连接失败, errMsg: %s", err.Error())
		panic(err)
	}

	err = conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		utils.Log().Panic("", "Kafka 设置读写超时时间失败, errMsg: %s", err.Error())
		panic(err)
	}

	KafkaConn = conn
}
