package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
)

func registerProductRouter(router *gin.RouterGroup, productHandler *handler.ProductHandler, jwtManager *jwt.JWTManager) {
	products := router.Group("/products")
	products.Use(middleware.JWTRegister(jwtManager))
	{
		products.GET("/:id", productHandler.Get)
		products.GET("", productHandler.GetAll)
		products.GET("/search", productHandler.Search)
		products.POST("", middleware.RequireRole("seller", "admin"), productHandler.Create)
		products.PUT("/:id", middleware.RequireRole("seller", "admin"), productHandler.Update)
		products.DELETE("/:id", middleware.RequireRole("seller", "admin"), productHandler.Delete)
	}
}
