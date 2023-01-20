package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	userRouter := r.Group("/users")
	{
		userRouter.POST("/", handlers.CreateUser)
		userRouter.POST("/login", handlers.Login)
		userRouter.GET("/", handlers.GetUser)
		userRouter.PUT("/", handlers.UpdateUser)
		userRouter.DELETE("/", handlers.DeleteUser)
	}
}
