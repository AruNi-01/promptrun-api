package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

func FindOrderListByBuyerId(c *gin.Context, buyerId int) ([]model.Order, *errs.Errs) {
	var orders []model.Order
	if model.DB.Where("buyer_id = ?", buyerId).Find(&orders).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取订单列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取订单列表失败"))
	}
	return orders, nil
}
