package service

import (
	"errors"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

func AddOrderRating(order model.Order) (bool, *errs.Errs) {
	if e := model.DB.Create(&model.OrderRating{
		OrderId:    order.Id,
		PromptId:   order.PromptId,
		SellerId:   order.SellerId,
		Rating:     order.Rating,
		CreateTime: time.Now(),
	}).Error; e != nil {
		utils.Log().Error("", "添加订单评分失败, err: %s", e.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("添加订单评分失败"))
	}
	return true, nil
}
