package cron

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"math"
	"promptrun-api/service"
	"promptrun-api/third_party/cache"
	"promptrun-api/utils"
	"strconv"
	"time"
)

const (
	RefreshPromptScoreCron = "0 */1 * * *" // 每小时执行
)

// PromptRun 纪元，用于计算 Prompt 的 score
var promptRunEpoch time.Time

func initEpoch(location *time.Location) {
	promptRunEpoch = time.Date(2001, time.October, 12, 0, 0, 0, 0, location)
}

// 低峰期使用定时任务访问缓存，获取期间发送变化的 Prompt，刷新 score
func refreshPromptScoreJob() {
	location := utils.GetShanghaiLocation()
	initEpoch(location)

	c := cron.New(cron.WithLocation(location))
	// 每小时执行
	_, err := c.AddFunc(RefreshPromptScoreCron, func() {
		utils.Log().Info("", "【定时任务】refreshPromptScoreJob start...")
		refreshPromptScore()
		utils.Log().Info("", "【定时任务】refreshPromptScoreJob end...")
	})
	if err != nil {
		utils.Log().Panic("", "【定时任务】Add refreshPromptScoreJob error: %s", err.Error())
		panic(err)
	}
	c.Start()
}

func refreshPromptScore() {
	// 从缓存中获取分数发送变化的 Prompt
	redisKey := cache.PromptScoreChangeKey()
	for {
		value, err := cache.RedisCli.SPop(context.Background(), redisKey).Result()
		if err != nil || value == "" {
			break
		}
		promptId, _ := strconv.Atoi(value)

		refresh(promptId)
	}
}

func refresh(promptId int) {
	prompt, err := service.FindPromptById(&gin.Context{}, promptId)
	if err != nil {
		utils.Log().Error("", "【定时任务】FindPromptById error: %s", err.Err.Error())
		return
	}

	rating := prompt.Rating
	browseAmount := prompt.BrowseAmount
	likeAmount := prompt.LikeAmount
	sellAmount := prompt.SellAmount
	createTime := prompt.CreateTime

	// 计算权重，权重 = 评分 * 10 + 销量 * 5 + 点赞数 * 3 + 浏览数 * 1
	weight := int(rating*10) + sellAmount*5 + likeAmount*3 + browseAmount*1
	if weight < 1 {
		weight = 1
	}

	// 计算 score = 权重（log 函数，权重越高增率越低） + 时效性（该 Prompt 创建时间距离 PromptRun 纪元的小时数，越新发布的时效性越高，score 也就相对越高）
	score := math.Log10(float64(weight)) + createTime.Sub(promptRunEpoch).Seconds()/3600/24

	// 更新 score
	if _, err := service.UpdatePromptScore(promptId, score); err != nil {
		utils.Log().Error("", "【定时任务】UpdatePromptScore error: %s", err.Err.Error())
		return
	}
}
