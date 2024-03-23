package routers

import (
	"github.com/gin-gonic/gin"
	"promptrun-api/api"
	"promptrun-api/middleware"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.Use(middleware.Cors()) // 跨域配置

	// 这种路由组方式方便直接通过搜索路由排查问题，用 Group 分开不利于查找定位
	rootGroup := route.Group("")
	{
		passport := rootGroup.Group("")
		{
			passport.POST("/api/v1/passport/register", api.Register)
			passport.POST("/api/v1/passport/login", api.Login)
		}

		prompt := rootGroup.Group("")
		{
			prompt.GET("/api/v1/prompt/list", api.PromptList)
		}

		// 需要登录拦截的路由
		auth := rootGroup.Group("")
		auth.Use(middleware.LoginRequired())
		{
			passport2 := auth.Group("")
			{
				passport2.GET("/api/v1/passport/logout", api.Logout)
			}
		}

	}
	return route
}
