package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
)

// LantuWxPay 蓝兔微信支付接口，返回支付二维码
func LantuWxPay(c *gin.Context) {
	var lantuWxPayReq service.LantuWxPayReq
	if err := c.ShouldBindJSON(&lantuWxPayReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		data, e := lantuWxPayReq.LantuWxPay(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(data))
	}
}

// LantuWxPayNotify 蓝兔微信支付回调通知
func LantuWxPayNotify(c *gin.Context) {
	var lantuWxPayNotifyParams service.LantuWxPayNotifyParams
	if err := c.ShouldBindJSON(&lantuWxPayNotifyParams); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, "FAIL")
	} else {
		_, e := lantuWxPayNotifyParams.LantuWxPayNotify(c)
		if e != nil {
			c.JSON(http.StatusOK, "FAIL")
			return
		}
		c.JSON(http.StatusOK, "SUCCESS")
	}
}
