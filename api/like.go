package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
	"strconv"
)

func IsLike(c *gin.Context) {
	userIdStr := c.Query("userId")
	promptIdStr := c.Query("promptId")
	// 参数为空则直接返回未点赞
	if userIdStr == "" || promptIdStr == "" {
		c.JSON(http.StatusOK, SuccessResponse(false))
		return
	}

	userId, _ := strconv.Atoi(userIdStr)
	promptId, _ := strconv.Atoi(promptIdStr)
	isLike, err := service.IsLike(c, userId, promptId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorResponse(err.ErrCode, err.Err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(isLike))
}

func Like(c *gin.Context) {
	var likeReq service.LikeReq
	if err := c.ShouldBindJSON(&likeReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		err := likeReq.Like(c)
		if err != nil {
			c.JSON(http.StatusOK, ErrorResponse(err.ErrCode, err.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}
