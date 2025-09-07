package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func registerCategoriesRouter(router *gin.RouterGroup, categoriesHandler *handler.CategoriesHandler, jwtManager *jwt.JWTManager, cache *redis.Client) {
	categories := router.Group("/categories")
	categories.Use(middleware.JWTRegister(jwtManager, cache))
	{
		categories.GET("", categoriesHandler.GetAll)
		admin := categories.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("", categoriesHandler.Create)
			admin.GET("/:id", categoriesHandler.GetById)
			admin.PUT("/:id", categoriesHandler.Update)
			admin.DELETE("/:id", categoriesHandler.Delete)
		}
	}
}
