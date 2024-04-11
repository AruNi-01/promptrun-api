package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
	"strconv"
)

func FindSellerByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	seller, e := service.FindSellerByUserId(c, userId)
	if e != nil {
		utils.Log().Error(c.FullPath(), e.Err.Error())
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(seller))
	}
}

func UpdateSeller(c *gin.Context) {
	var sellerUpdateReq service.SellerUpdateReq
	if err := c.ShouldBindJSON(&sellerUpdateReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		_, e := sellerUpdateReq.UpdateSeller(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

func FindSellerById(c *gin.Context) {
	sellerId, _ := strconv.Atoi(c.Param("id"))
	seller, e := service.FindSellerById(c, sellerId)
	if e != nil {
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(seller))
}
