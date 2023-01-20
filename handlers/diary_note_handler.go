package handlers

import (
	// "remood/models"

	"github.com/gin-gonic/gin"
)

func CreateDiaryNote(ctx *gin.Context) {
}

func GetAllDiaryNotes(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Hello from Get All Diary Notes"})
}

func GetDiaryNote(ctx *gin.Context) {

}

func GetSomeDiaryNote(ctx *gin.Context) {

}

func UpdateDiaryNote(ctx *gin.Context) {

}

func DeleteDiaryNote(ctx *gin.Context) {

}
