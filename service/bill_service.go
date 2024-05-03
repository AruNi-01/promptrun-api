package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

func AddBill(c *gin.Context, bill model.Bill) (bool, *errs.Errs) {
	if err := model.DB.Create(&bill).Error; err != nil {
		utils.Log().Error(c.FullPath(), "创建账单失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("创建账单失败"))
	}
	return true, nil
}
