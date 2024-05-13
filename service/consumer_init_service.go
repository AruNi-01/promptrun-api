package service

import (
	"encoding/json"
	"promptrun-api/common/constants"
	"promptrun-api/third_party/kafka2"
	"promptrun-api/third_party/kafka2/vos"
	"promptrun-api/utils"
)

func InitConsumer() {
	promptBuyerRatingConsumerInit()
}

func promptBuyerRatingConsumerInit() {
	kafka2.Subscribe(constants.PromptBuyerRatingTopic, func(msg string) {
		var buyerRatingResult vos.BuyerRatingResult
		if json.Unmarshal([]byte(msg), &buyerRatingResult) != nil {
			utils.Log().Error("", "【MQ 消费】, 反序列化消息失败")
			return
		}

		order := buyerRatingResult.Order
		// 插入订单评分表
		if _, e := AddOrderRating(order); e != nil {
			utils.Log().Error("", "【MQ 消费】, 异步插入订单评分表失败, errMsg: %s", e.Err.Error())
		}

		// 插入评分消息
		errs := OrderRatingMsgNotice(order)
		if errs != nil {
			utils.Log().Error("", "【MQ 消费】, 异步发送评分消息失败, errMsg: %s", errs.Err.Error())
		}
	})
}
