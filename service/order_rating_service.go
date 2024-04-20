package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

func AddOrderRating(c *gin.Context, order model.Order, rating float64) (bool, *errs.Errs) {
	if e := model.DB.Create(&model.OrderRating{
		OrderId:    order.Id,
		PromptId:   order.PromptId,
		SellerId:   order.SellerId,
		Rating:     rating,
		CreateTime: time.Now(),
	}).Error; e != nil {
		utils.Log().Error(c.FullPath(), "添加订单评分失败, err: %s", e.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("添加订单评分失败"))
	}
	return true, nil
}
