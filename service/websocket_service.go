package service

import (
	"github.com/gin-gonic/gin"
	"promptrun-api/utils/websocket2"
)

const (
	MessageNotReadCountWsPrefix = "ws_message_not_read_count_"
)

func WsMessageNotReadCountByUserId(c *gin.Context, userId string) {
	websocket2.WsHandler(c.Writer, c.Request, MessageNotReadCountWsPrefix+userId)
}
