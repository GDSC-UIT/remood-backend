package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func UserRouter(gin *gin.RouterGroup) {
	userRouter := gin.Group("/user")
	{
		userRouter.POST("/", handlers.CreateUser)
	}
}
