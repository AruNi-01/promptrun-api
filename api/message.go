package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/service"
	"promptrun-api/utils"
	"strconv"
)

func MessageListByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	messageList, e := service.FindMessageListByUserId(c, userId)
	if e != nil {
		utils.Log().Error(c.FullPath(), e.Err.Error())
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(messageList))
	}
}

func ReadAllMessage(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	if e := service.ReadAllMessage(c, userId); e != nil {
		utils.Log().Error(c.FullPath(), e.Err.Error())
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

func MessageNotReadCountByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	count, e := service.FindMessageNotReadListByUserId(c, userId)
	if e != nil {
		utils.Log().Error(c.FullPath(), e.Err.Error())
		c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(count))
	}
}
