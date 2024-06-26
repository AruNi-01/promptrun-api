package routers

import (
	"github.com/gin-gonic/gin"
	"os"
	"promptrun-api/api"
	"promptrun-api/middleware"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))

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

		user := rootGroup.Group("")
		{
			user.GET("api/v1/user/findById/:id", api.FindUserById)
		}

		seller := rootGroup.Group("")
		{
			seller.GET("api/v1/seller/findById/:id", api.FindSellerById)
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
			prompt.POST("/api/v1/prompt/listBySellerId", api.FindListBySellerId)

			prompt.POST("/api/v1/prompt/updateBrowseAmountById/:id", api.UpdateBrowseAmountById)
		}

		model := rootGroup.Group("")
		{
			model.GET("api/v1/model/list", api.ModelList)
			model.GET("api/v1/model/findById/:id", api.FindModelById)
		}

		like := rootGroup.Group("")
		{
			like.GET("api/v1/likes/isLike", api.IsLike)
		}

		order := rootGroup.Group("")
		{
			order.GET("api/v1/order/findById/:orderId", api.FindOrderById)
		}

		pay := rootGroup.Group("")
		{
			pay.GET("api/v1/pay/lantuWxPayNotify", api.LantuWxPayNotify)
		}

		// 需要登录拦截的路由
		auth := rootGroup.Group("")
		auth.Use(middleware.LoginRequired())
		{
			passport2 := auth.Group("")
			{
				passport2.GET("/api/v1/passport/logout", api.Logout)
				passport2.POST("/api/v1/passport/updatePassword", api.UpdatePassword)

				passport2.GET("api/v1/passport/checkIsLogin", api.CheckIsLogin)
			}

			user2 := auth.Group("")
			{
				user2.POST("api/v1/user/update", api.UpdateUser)
				user2.POST("api/v1/user/becomeSeller", api.UserBecomeSeller)
			}

			seller2 := auth.Group("")
			{
				seller2.GET("api/v1/seller/findByUserId/:userId", api.FindSellerByUserId)
				seller2.POST("api/v1/seller/update", api.UpdateSeller)
			}

			order2 := auth.Group("")
			{
				order2.POST("api/v1/order/listAttachFullInfoBySellerUserId", api.OrderListAttachFullInfoBySellerUserId)
				order2.GET("api/v1/order/findChartsFullInfoBySellerUserId/:sellerUserId", api.FindChartsFullInfoBySellerUserId)
				order2.GET("api/v1/order/listAttachPromptDetailById/:orderId", api.OrderListAttachPromptDetailById)
				order2.POST("api/v1/order/ratingById", api.OrderRatingById)
				order2.GET("api/v1/order/isRepeatPurchase", api.IsRepeatPurchase)
			}

			like2 := auth.Group("")
			{
				like2.POST("api/v1/likes/like", api.Like)
			}

			prompt2 := auth.Group("")
			{
				prompt2.POST("api/v1/prompt/publish", api.PromptPublish)
				prompt2.POST("api/v1/prompt/updatePublishStatusById", api.UpdatePublishStatusById)
			}

			pay2 := rootGroup.Group("")
			{
				pay2.POST("api/v1/pay/lantuWxPay", api.LantuWxPay)
				pay2.POST("api/v1/pay/lantuWxPayQueryOrder", api.LantuWxPayQueryOrder)
				pay2.POST("api/v1/pay/balancePay", api.BalancePay)
			}

			wallet2 := rootGroup.Group("")
			{
				wallet2.GET("api/v1/wallet/findByUserId/:userId", api.FindByUserId)
			}

			bill2 := rootGroup.Group("")
			{
				bill2.POST("api/v1/bill/listByUserId", api.BillListByUserId)
			}

			message2 := rootGroup.Group("")
			{
				message2.GET("api/v1/message/listByUserId/:userId", api.MessageListByUserId)
				message2.POST("api/v1/message/readAll/:userId", api.ReadAllMessage)
				message2.GET("api/v1/message/notReadCountByUserId/:userId", api.MessageNotReadCountByUserId)
			}

			websocket2 := rootGroup.Group("")
			{
				websocket2.GET("api/v1/websocket/messageNotReadCountByUserId/:userId", api.WsMessageNotReadCountByUserId)
			}
		}

	}
	return route
}
