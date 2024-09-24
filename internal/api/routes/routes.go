package routes

import (
	"github.com/bimaputraas/rest-api/internal/api/controller"
	"github.com/bimaputraas/rest-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func New(middleware *middleware.Middleware, controller *controller.Controller) *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(middleware.Cors())

	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	auth := router.Group("").Use(middleware.Auth())
	auth.POST("/topup", controller.TopUp)
	auth.POST("/pay", controller.Payment)
	auth.POST("/transfer", controller.Transfer)
	auth.GET("/transactions", controller.GetUserTransactions)

	return router
}
