package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/cache"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	"promptrun-api/utils"
)

func Register(c *gin.Context) {
	var registerService service.RegisterReq
	if err := c.ShouldBindJSON(&registerService); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		user, e := registerService.Register(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(user))
	}
}

func Login(c *gin.Context) {
	var loginService service.LoginReq
	if err := c.ShouldBindJSON(&loginService); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		user, e := loginService.Login(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(user))
	}
}

func Logout(c *gin.Context) {
	// 直接删除 Redis 中的登录凭证
	ticket, err := c.Cookie("ticket")
	if err != nil {
		c.JSON(http.StatusOK, ErrorResponse(errs.ErrNotLogin, "您未登录，请登录后再操作"))
		return
	}
	ticketKey := cache.TicketKey(ticket)
	cache.RedisCli.Del(c, ticketKey)

	c.JSON(http.StatusOK, SuccessResponse(nil))
}
