package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
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
