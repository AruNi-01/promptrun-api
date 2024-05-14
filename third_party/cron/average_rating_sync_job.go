package cron

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"promptrun-api/service"
	"promptrun-api/utils"
)

const (
	SyncAverageRatingCron = "0 2 * * *" // 每天凌晨 2 点
)

// 低峰期使用定时任务扫描 order_rating 表，同步评分到 seller、prompt 表（计算平均评分）
func syncAverageRatingJob() {
	location := utils.GetShanghaiLocation()
	c := cron.New(cron.WithLocation(location))
	// 每天凌晨 2 点执行
	_, err := c.AddFunc(SyncAverageRatingCron, func() {
		utils.Log().Info("", "【定时任务】syncAverageRatingJob start...")
		syncAverageRating()
	})
	if err != nil {
		utils.Log().Panic("", "【定时任务】Add syncAverageRatingJob error: %s", err.Error())
		panic(err)
	}
	c.Start()
}

func syncAverageRating() {
	// 扫描 order_rating 表，获取昨天的订单评分
	orderRatings, errs := service.GetOrderRatingYesterday()
	if errs != nil {
		utils.Log().Error("", "【定时任务】GetOrderRatingYesterday error: %s", errs.Err.Error())
		return
	}

	// 按卖家和 Prompt 分组
	sellerRatingMap := make(map[int][]float64)
	promptRatingMap := make(map[int][]float64)
	for _, orderRating := range orderRatings {
		if _, ok := sellerRatingMap[orderRating.SellerId]; !ok {
			sellerRatingMap[orderRating.SellerId] = make([]float64, 0)
		}
		sellerRatingMap[orderRating.SellerId] = append(sellerRatingMap[orderRating.SellerId], orderRating.Rating)

		if _, ok := promptRatingMap[orderRating.PromptId]; !ok {
			promptRatingMap[orderRating.PromptId] = make([]float64, 0)
		}
		promptRatingMap[orderRating.PromptId] = append(promptRatingMap[orderRating.PromptId], orderRating.Rating)
	}

	go calculateSellerRating(sellerRatingMap)
	go calculatePromptRating(promptRatingMap)
}

// 计算 seller 平均评分
func calculateSellerRating(sellerRatingMap map[int][]float64) {
	for sellerId, ratings := range sellerRatingMap {
		averageRating := 0.0
		for _, rating := range ratings {
			averageRating += rating
		}
		averageRating /= float64(len(ratings))

		// 查询出之前的评分，计算平均评分
		seller, errs := service.FindSellerById(&gin.Context{}, sellerId)
		if errs != nil {
			utils.Log().Error("", "【定时任务】FindSellerById error: %s", errs.Err.Error())
			continue
		}
		averageRating = (averageRating + seller.Rating) / 2
		_, e := service.UpdateSellerRating(sellerId, averageRating)
		if e != nil {
			utils.Log().Error("", "【定时任务】UpdateSellerRating error: %s", e.Err.Error())
			continue
		}
	}
}

// 计算 Prompt 平均评分
func calculatePromptRating(promptRatingMap map[int][]float64) {
	for promptId, ratings := range promptRatingMap {
		averageRating := 0.0
		for _, rating := range ratings {
			averageRating += rating
		}
		averageRating /= float64(len(ratings))

		// 查询出之前的评分，计算平均评分
		prompt, errs := service.FindPromptById(&gin.Context{}, promptId)
		if errs != nil {
			utils.Log().Error("", "【定时任务】FindPromptById error: %s", errs.Err.Error())
			continue
		}
		averageRating = (averageRating + prompt.Rating) / 2
		_, e := service.UpdatePromptRating(promptId, averageRating)
		if e != nil {
			utils.Log().Error("", "【定时任务】UpdatePromptRating error: %s", e.Err.Error())
			continue
		}
	}
}
