package handlers

import (
	// "fmt"
	"remood/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context) {
	var user models.User
	if ctx.ShouldBindJSON(&user) != nil {
		return
	}
	err := user.Create()
	if err != nil {
		return
	}

	ctx.JSON(200, models.Response{
		Message: "Create User Successfully",
		Error:   false,
		Data:    user,
	})
}

func Login(ctx *gin.Context) {
}

func GetUser(ctx *gin.Context) {
}

func UpdateUser(ctx *gin.Context) {
}

func DeleteUser(ctx *gin.Context) {
}
