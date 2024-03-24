package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

func ModelList(c *gin.Context) ([]model.Model, *errs.Errs) {
	var models []model.Model
	if model.DB.Find(&models).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取模型列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取模型列表失败"))
	}
	return models, nil
}

func FindModelById(id int) (model.Model, *errs.Errs) {
	var promptModel model.Model
	if model.DB.First(&promptModel, id).RecordNotFound() {
		return model.Model{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该模型"))
	}
	return promptModel, nil
}
