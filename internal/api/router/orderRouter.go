package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
)

func registerOrderRouter(router *gin.RouterGroup, orderHandler *handler.OrderHandler, jwtManager *jwt.JWTManager) {
	order := router.Group("/order")
	order.Use(middleware.JWTRegister(jwtManager))
	{
		order.POST("", orderHandler.Create)
		order.GET("/history", orderHandler.GetOrdersByUserId)
		order.GET("/items/:id", orderHandler.GetOrderItemsByOrderId)
		order.GET("/:id", orderHandler.GetOrderById)
		order.DELETE("/:id", orderHandler.DeleteOrderById)
	}
}
