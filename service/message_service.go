package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

const (
	MessageFromSystem = 0
)

func FindMessageListByUserId(c *gin.Context, userId int) ([]model.Message, *errs.Errs) {
	var messageList []model.Message
	if err := model.DB.Where("to_user_id = ?", userId).
		Limit(200).
		Order("create_time ASC").
		Find(&messageList).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 查询消息列表失败，errMsg: %s", err.Error())
		return nil, errs.NewErrs(errs.ErrDBError, err)
	}
	return messageList, nil
}

func ReadAllMessage(c *gin.Context, userId int) *errs.Errs {
	if err := model.DB.Model(&model.Message{}).
		Where("to_user_id = ?", userId).
		Update("is_read", model.MessageIsReadYes).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 更新消息已读失败，errMsg: %s", err.Error())
		return errs.NewErrs(errs.ErrDBError, err)
	}
	return nil
}

func FindMessageNotReadListByUserId(c *gin.Context, userId int) (int, *errs.Errs) {
	var count int
	if err := model.DB.Model(&model.Message{}).
		Where("to_user_id = ? AND is_read = ?", userId, model.MessageIsReadNo).
		Count(&count).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 查询未读消息数量失败，errMsg: %s", err.Error())
		return 0, errs.NewErrs(errs.ErrDBError, err)
	}
	return count, nil
}

func RegisterMsgNotice(c *gin.Context, user model.User) {
	message := model.Message{
		FromUserId: MessageFromSystem,
		ToUserId:   user.Id,
		Type:       model.MessageTypeActivity,
		Content:    fmt.Sprintf("Hello %s，欢迎加入 PromptRun，赠送的 ￥10.00 已放入我的钱包-余额中，快去交易市场挑选喜欢的 Prompt 吧！", user.Nickname),
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "注册 Prompt 通知链路发生错误 -> 创建消息失败")
		return
	}
}

func SellerBecomeMsgNotice(c *gin.Context, userId int) {
	message := model.Message{
		FromUserId: MessageFromSystem,
		ToUserId:   userId,
		Type:       model.MessageTypeSellerBecome,
		Content:    "恭喜您，已成功成为卖家，快去发布您的 Prompt 吧！",
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "成为卖家 Prompt 通知链路发生错误 -> 创建消息失败")
		return
	}

}

func LikeMsgNotice(c *gin.Context, likes model.Likes) {
	seller, e := FindSellerById(c, likes.SellerId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "点赞 Prompt 通知链路发生错误 -> 未找到卖家")
		return
	}

	user, e := FindUserById(c, likes.UserId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "点赞 Prompt 通知链路发生错误 -> 未找到用户")
		return
	}

	prompt, e := FindPromptById(c, likes.PromptId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "点赞 Prompt 通知链路发生错误 -> 未找到 Prompt")
		return
	}

	message := model.Message{
		FromUserId: MessageFromSystem,
		ToUserId:   seller.UserId,
		Type:       model.MessageTypeLike,
		Content:    fmt.Sprintf("%s 喜欢了您的 Prompt - %s！", user.Nickname, prompt.Title),
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "点赞 Prompt 通知链路发生错误 -> 创建消息失败")
		return
	}
}

func OrderRatingMsgNotice(c *gin.Context, order model.Order) {
	seller, e := FindSellerById(c, order.SellerId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "评价订单通知链路发生错误 -> 未找到卖家")
		return
	}

	user, e := FindUserById(c, order.BuyerId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "评价订单通知链路发生错误 -> 未找到用户")
		return
	}

	prompt, e := FindPromptById(c, order.PromptId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "评价订单通知链路发生错误 -> 未找到 Prompt")
		return
	}

	message := model.Message{
		FromUserId: MessageFromSystem,
		ToUserId:   seller.UserId,
		Type:       model.MessageTypeOrderRating,
		Content:    fmt.Sprintf("%s 评价了您的 Prompt - %s，评分：%.1f！", user.Nickname, prompt.Title, order.Rating),
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "评价订单通知链路发生错误 -> 创建消息失败")
		return
	}
}

func PromptSoldMsgNotice(c *gin.Context, promptTitle string, sellerUserId, buyerId int) {
	user, e := FindUserById(c, buyerId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "Prompt 售出通知链路发生错误 -> 未找到用户")
		return
	}

	message := model.Message{
		FromUserId: MessageFromSystem,
		ToUserId:   sellerUserId,
		Type:       model.MessageTypeSell,
		Content:    fmt.Sprintf("%s 购买了您的 Prompt - %s，获得的收入已放入我的钱包-余额中！", user.Nickname, promptTitle),
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "Prompt 售出通知链路发生错误 -> 创建消息失败")
		return
	}
}

func PromptPublishMsgNotice(c *gin.Context, promptId int) {
	prompt, e := FindPromptById(c, promptId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "Prompt 发布通知链路发生错误 -> 未找到 Prompt")
		return
	}

	seller, e := FindSellerById(c, prompt.SellerId)
	if e != nil {
		utils.Log().Error(c.FullPath(), "Prompt 发布通知链路发生错误 -> 未找到卖家")
		return

	}

	message := model.Message{
		FromUserId: MessageFromSystem,
		ToUserId:   seller.UserId,
		Type:       model.MessageTypePublish,
		Content:    fmt.Sprintf("您的 Prompt - %s 发布成功，已在交易市场开始售卖中！", prompt.Title),
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "Prompt 发布通知链路发生错误 -> 创建消息失败")
		return
	}
}
