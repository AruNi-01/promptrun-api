package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
