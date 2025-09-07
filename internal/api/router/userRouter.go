package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/api/middleware"
	"github.com/niklvrr/myMarketplace/internal/handler"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
)

func registerUserRouter(router *gin.RouterGroup, userHandler *handler.UserHandler, jwtManager *jwt.JWTManager) {
	user := router.Group("/user")
	{
		user.POST("/signup", userHandler.SignUp)
		user.POST("/login", userHandler.Login)

		auth := user.Group("")
		auth.Use(middleware.JWTRegister(jwtManager))
		{
			auth.GET("", userHandler.GetUserById)
			auth.PUT("", userHandler.UpdateUserById)
			auth.POST("", userHandler.GetUserByEmail)
			auth.PUT("/role", userHandler.UpdateUserRole)
			auth.POST("/logout", userHandler.Logout)

			admin := auth.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.PUT("/block", userHandler.BlockUserById)
				admin.PUT("/unblock", userHandler.UnblockUserById)
				admin.GET("", userHandler.GetAllUsers)
				admin.PUT("/approve", userHandler.ApproveProduct)
			}
		}
	}
}
