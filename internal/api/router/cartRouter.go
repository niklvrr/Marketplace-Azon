package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
)

func registerCartRouter(router *gin.RouterGroup, cartHandler *handler.CartHandler, jwtManager *jwt.JWTManager) {
	cart := router.Group("/cart")
	cart.Use(middleware.JWTRegister(jwtManager))
	{
		cart.GET("", cartHandler.GetCartByUserId)
		cart.GET("/:id", cartHandler.GetCartItemsByCartId)
		cart.POST("", cartHandler.AddItem)
		cart.DELETE("", cartHandler.RemoveItem)
		cart.DELETE("/clear", cartHandler.ClearCart)
	}
}
