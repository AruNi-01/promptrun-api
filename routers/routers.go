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
			prompt.POST("/api/v1/prompt/list", api.PromptList)
			prompt.GET("/api/v1/prompt/findById/:id", api.FindById)
			prompt.GET("/api/v1/prompt/findFullInfoById/:id", api.FindFullInfoById)

			prompt.GET("/api/v1/prompt/findImgListByPromptId/:id", api.FindImgListByPromptId)
			prompt.GET("/api/v1/prompt/findMasterImgByPromptId/:id", api.FindMasterImgByPromptId)
			prompt.POST("/api/v1/prompt/findMasterImgListByPromptIds", api.FindMasterImgListByPromptIds)

			prompt.POST("/api/v1/prompt/listByBuyerId", api.FindListByBuyerId)
		}

		model := rootGroup.Group("")
		{
			model.GET("api/v1/model/list", api.ModelList)
			model.GET("api/v1/model/findById/:id", api.FindModelById)
		}

		// 需要登录拦截的路由
		auth := rootGroup.Group("")
		auth.Use(middleware.LoginRequired())
		{
			passport2 := auth.Group("")
			{
				passport2.GET("/api/v1/passport/logout", api.Logout)
			}

			user2 := auth.Group("")
			{
				user2.GET("api/v1/user/findById/:id", api.FindUserById)
				user2.POST("api/v1/user/update", api.UpdateUser)
			}
		}

	}
	return route
}
