package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/constants"
	"promptrun-api/model"
	"promptrun-api/third_party/kafka2"
	"promptrun-api/third_party/kafka2/vos"
	"promptrun-api/utils"
	"time"
)

func InitConsumer() {
	promptBuyerRatingConsumerInit()
	promptSoldResultConsumerInit()
}

func promptSoldResultConsumerInit() {
	kafka2.Subscribe(constants.PromptSoldResultTopic, func(msg string) {
		var soldResult vos.PromptSoldResult
		if json.Unmarshal([]byte(msg), &soldResult) != nil {
			utils.Log().Error("", "【MQ 消费】, 反序列化消息失败")
			return
		}

		// 增加 Prompt 和卖家销售量
		IncreasePromptSellAmount(soldResult.PromptId)
		IncreaseSellerSellAmount(soldResult.SellerId)

		// 生成账单
		bill := model.Bill{
			UserId:     soldResult.SellerUserId,
			Type:       model.BillTypeIncome,
			Amount:     soldResult.Price,
			Channel:    soldResult.IncomeChannel,
			Remark:     fmt.Sprintf("售出 Prompt - %s", soldResult.PromptTitle),
			CreateTime: time.Now(),
		}
		if _, e := AddBill(&gin.Context{}, bill); e != nil {
			utils.Log().Error("", "【MQ 消费】Add seller bill error: %s", e.Err.Error())
		}

		bill.UserId = soldResult.BuyerId
		bill.Type = model.BillTypeOutcome
		bill.Channel = soldResult.OutcomeChannel
		bill.Remark = fmt.Sprintf("购买 Prompt - %s", soldResult.PromptTitle)
		if _, e := AddBill(&gin.Context{}, bill); e != nil {
			utils.Log().Error("", "【MQ 消费】Add buyer bill error: %s", e.Err.Error())
		}

		// 插入售出消息
		PromptSoldMsgNotice(soldResult.PromptTitle, soldResult.SellerUserId, soldResult.BuyerId)
	})
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
