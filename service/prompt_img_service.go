package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

func FindPromptImgListByPromptId(c *gin.Context, id int) ([]model.PromptImg, *errs.Errs) {
	var promptImgList []model.PromptImg
	if model.DB.Where("prompt_id = ?", id).Find(&promptImgList).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示词图片列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词图片列表失败"))
	}
	return promptImgList, nil
}
