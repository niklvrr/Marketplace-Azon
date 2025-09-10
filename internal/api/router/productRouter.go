package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler/productHandler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func registerProductRouter(router *gin.RouterGroup, productHandler *productHandler.ProductHandler, jwtManager *jwt.JWTManager, cache *redis.Client) {
	products := router.Group("/products")
	products.Use(middleware.JWTRegister(jwtManager, cache))
	{
		products.GET("/:id", productHandler.Get)
		products.GET("", productHandler.GetAll)
		products.GET("/search", productHandler.Search)

		seller := products.Group("/")
		seller.Use(middleware.RequireRole("seller", "admin"))
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}
}
