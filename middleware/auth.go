package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/result"
	"promptrun-api/service"
	"time"
)

// LoginRequired 登录校验中间件
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Cookie 中的 ticket（登录凭证）
		ticket, err := c.Cookie("ticket")
		if err != nil {
			c.JSON(http.StatusOK, result.NotLogin())
			c.Abort()
			return
		}

		// 根据 ticket 从 Redis 中获取 LoginTicket
		loginTicket, e := service.FindLoginTicket(c, ticket)
		if e != nil { // Json unmarshal error
			c.JSON(http.StatusOK, result.Err(e.ErrCode, e.Err.Error()))
			c.Abort()
			return
		} else if loginTicket.Ticket == "" || loginTicket.ExpiredAt.Before(time.Now()) {
			c.JSON(http.StatusOK, result.NotLogin())
			c.Abort()
			return
		}
		c.Next()
	}
}
