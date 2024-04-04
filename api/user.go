package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
	"strconv"
)

func FindUserById(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	user, e := service.FindUserById(c, userId)
	if e != nil {
		utils.Log().Error(c.FullPath(), e.Err.Error())
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(user))
	}
}

func UpdateUser(c *gin.Context) {
	var userUpdateReq service.UserUpdateReq
	if err := c.ShouldBindJSON(&userUpdateReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		_, e := userUpdateReq.UpdateUser(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}
