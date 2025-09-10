package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler/cartHandler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func registerCartRouter(router *gin.RouterGroup, cartHandler *cartHandler.CartHandler, jwtManager *jwt.JWTManager, cache *redis.Client) {
	cart := router.Group("/cart")
	cart.Use(middleware.JWTRegister(jwtManager, cache))
	{
		cart.GET("", cartHandler.GetCartByUserId)
		cart.GET("/:id", cartHandler.GetCartItemsByCartId)
		cart.POST("", cartHandler.AddItem)
		cart.DELETE("", cartHandler.RemoveItem)
		cart.DELETE("/clear", cartHandler.ClearCart)
	}
}
