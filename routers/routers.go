package routers

import (
	"github.com/gin-gonic/gin"
	"promptrun-api/api"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	// 这种路由组方式方便直接通过搜索路由排查问题，用 Group 分开不利于查找定位
	rootGroup := route.Group("")
	{
		user := rootGroup.Group("")
		{
			user.POST("/api/v1/user/register", api.Register)
			user.POST("/api/v1/user/login", api.Login)
		}
	}
	return route
}
