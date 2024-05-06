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
		Content:    fmt.Sprintf("%s 喜欢了您的 Prompt - %s", user.Nickname, prompt.Title),
		IsRead:     model.MessageIsReadNo,
		CreateTime: time.Now(),
	}

	if err := model.DB.Create(&message).Error; err != nil {
		utils.Log().Error(c.FullPath(), "点赞 Prompt 通知链路发生错误 -> 创建消息失败")
		return
	}
}
