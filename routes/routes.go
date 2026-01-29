package routes

import (
	"payout-backend/controllers"
	"payout-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	public := r.Group("api")
	public.GET("/Hi", controllers.Hello)

	public.POST("login", controllers.Login)

	auth := r.Group("api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("dashboardData", controllers.DashboardData)
		auth.GET("merchants", controllers.GetMerchants)
		auth.GET("merchants/:id", controllers.GetMerchantsByID)
		auth.PUT("merchants/:id", controllers.UpdateMerchant)
		auth.POST("merchants", controllers.AddMerchant)
		auth.GET("merchantAccountDetails", controllers.GetMerchantAccountDetails)

		auth.GET("payoutProviders", controllers.GetPayoutProvider)
		auth.GET("payoutProviders/:id", controllers.GetPayoutProviderByID)
		auth.PUT("payoutProviders/:id", controllers.UpdatePayoutProvider)
		auth.POST("payoutProviders", controllers.AddPayoutProvider)

		auth.GET("providerType", controllers.GetProviderType)
		auth.GET("providerType/:id", controllers.GetProviderTypeByID)
		auth.PUT("providerType/:id", controllers.UpdateProviderType)
		auth.POST("providerType", controllers.AddProviderType)

		auth.GET("transaction", controllers.GetTransaction)
		auth.GET("transaction/:id", controllers.GetTransactionByID)

		auth.POST("logout", controllers.Logout)

	}

	payout := r.Group("api/payout")
	payout.Use(middleware.MerchantValidator())
	{
		payout.POST("make-payout", controllers.MakePayout)
	}

}
