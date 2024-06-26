package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"promptrun-api/common/errs"
	"promptrun-api/service"
	cache2 "promptrun-api/third_party/cache"
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
	ticketKey := cache2.TicketKey(ticket)
	cache2.RedisCli.Del(c, ticketKey)

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func UpdatePassword(c *gin.Context) {
	var updatePasswordReq service.UpdatePasswordReq
	if err := c.ShouldBindJSON(&updatePasswordReq); err != nil {
		utils.Log().Error(c.FullPath(), "请求参数错误")
		c.JSON(http.StatusBadRequest, ErrorResponse(errs.ErrParam, "请求参数错误"))
	} else {
		flag, e := updatePasswordReq.UpdatePassword(c)
		if e != nil {
			c.JSON(http.StatusOK, ErrorResponse(e.ErrCode, e.Err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(flag))
	}
}

// CheckIsLogin 检查是否登录，Auth 中间件会拦截未登录的请求
func CheckIsLogin(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse(nil))
}
