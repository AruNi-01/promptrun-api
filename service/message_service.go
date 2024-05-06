package service

import (
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
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
