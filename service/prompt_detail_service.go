package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

func FindPromptDetailByPromptId(c *gin.Context, promptId int) (model.PromptDetail, *errs.Errs) {
	var promptDetail model.PromptDetail
	if err := model.DB.Where("prompt_id = ?", promptId).First(&promptDetail).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示详情失败")
		return model.PromptDetail{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示详情失败"))
	}
	return promptDetail, nil
}
