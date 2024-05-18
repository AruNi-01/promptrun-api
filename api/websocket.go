package api

import (
	"github.com/gin-gonic/gin"
	"promptrun-api/service"
)

func WsMessageNotReadCountByUserId(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "0" {
		return
	}
	service.WsMessageNotReadCountByUserId(c, userId)
}
