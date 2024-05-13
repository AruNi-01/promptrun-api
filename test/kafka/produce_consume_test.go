package kafka

import (
	"encoding/json"
	"promptrun-api/common/constants"
	"promptrun-api/third_party/kafka2"
	"testing"
)

type SoldResult struct {
	OrderId string `json:"order_id"`
	Amount  int    `json:"amount"`
	Mark    string `json:"mark"`
}

func TestProduce(t *testing.T) {
	kafka2.InitKafkaWriter()

	result := SoldResult{
		OrderId: "100000000001",
		Amount:  10,
		Mark:    "test",
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		t.Errorf("json marshal error: %v", err)
	}
	if kafka2.SendMessage(constants.PromptSoldResultTopic, "", string(jsonResult)) != nil {
		t.Errorf("send message fail")
	}
}

func TestConsume(t *testing.T) {
	kafka2.Subscribe(constants.PromptSoldResultTopic, func(msg string) {
		t.Logf("received message: %s", msg)
	})
}
