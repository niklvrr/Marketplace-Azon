package router

import (
	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/handler"
)

func registerUserRouter(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	user := router.Group("/user")
	{
		user.POST("", userHandler.SignUp)
		user.GET("", userHandler.Login)
		user.GET("/:id", userHandler.GetUserById)
		user.PUT("/:id", userHandler.UpdateUserById)

		admin := user.Group("/admin")
		{
			admin.PUT("/block/:id", userHandler.BlockUserById)
			admin.PUT("/unblock/:id", userHandler.UnblockUserById)
			admin.GET("", userHandler.GetAllUsers)
			admin.PUT("/role/:id", userHandler.UpdateUserRole)
			admin.PUT("/approve/:id", userHandler.ApproveProduct)
		}
	}
}
