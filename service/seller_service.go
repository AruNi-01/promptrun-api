package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

func FindSellerById(c *gin.Context, id int) (model.Seller, *errs.Errs) {
	var seller model.Seller
	if model.DB.First(&seller, id).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "未找到该卖家")
		return model.Seller{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该卖家"))
	}
	return seller, nil
}

func FindSellerByUserId(c *gin.Context, userId int) (model.Seller, *errs.Errs) {
	var seller model.Seller
	if model.DB.Where("user_id = ?", userId).First(&seller).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "未找到该卖家")
		return model.Seller{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该卖家"))
	}
	return seller, nil
}
