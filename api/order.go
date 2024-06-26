package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
	"strconv"
)

func OrderListAttachFullInfoBySellerUserId(c *gin.Context) {
	var orderListBySellerUserIdReq service.OrderListBySellerUserIdReq
	if err := c.ShouldBindJSON(&orderListBySellerUserIdReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		orderListAttachFullInfo, e := orderListBySellerUserIdReq.FindOrderListAttachFullInfoBySellerUserId(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(orderListAttachFullInfo))
	}
}

func FindChartsFullInfoBySellerUserId(c *gin.Context) {
	sellerUserId, _ := strconv.Atoi(c.Param("sellerUserId"))
	chartsRsp, e := service.FindChartsFullInfoBySellerUserId(c, sellerUserId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, SuccessResponse(chartsRsp))
	}
}

func OrderListAttachPromptDetailById(c *gin.Context) {
	orderId, _ := strconv.ParseInt(c.Param("orderId"), 10, 64)
	orderAttachPromptDetail, e := service.FindOrderListAttachPromptDetailById(c, orderId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, SuccessResponse(orderAttachPromptDetail))
	}
}

func OrderRatingById(c *gin.Context) {
	orderId, _ := strconv.ParseInt(c.Query("orderId"), 10, 64)
	rating, _ := strconv.ParseFloat(c.Query("rating"), 64)
	order, e := service.OrderRatingById(c, orderId, rating)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, SuccessResponse(order))
	}
}

func FindOrderById(c *gin.Context) {
	orderId, _ := strconv.ParseInt(c.Param("orderId"), 10, 64)
	order, e := service.FindOrderById(c, orderId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, SuccessResponse(order))
	}
}

func IsRepeatPurchase(c *gin.Context) {
	buyerId, _ := strconv.Atoi(c.Query("buyerId"))
	promptId, _ := strconv.Atoi(c.Query("promptId"))
	isRepeatPurchase, err := service.IsRepeatPurchase(c, buyerId, promptId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorResponse(err.ErrCode, err.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(isRepeatPurchase))
}
