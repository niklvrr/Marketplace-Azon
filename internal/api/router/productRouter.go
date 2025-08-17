package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/handler"
)

func registerProductRouter(router *gin.RouterGroup, productHandler *handler.ProductHandler) {
	products := router.Group("/products")
	{
		products.GET("/:id", productHandler.Get)
		products.GET("", productHandler.GetAll)
		products.GET("/search", productHandler.Search)
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}
}
