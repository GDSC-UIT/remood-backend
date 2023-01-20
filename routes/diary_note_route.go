package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func DiaryNoteRouter(r *gin.RouterGroup) {
	diaryNoteRouter := r.Group("diary-notes")
	{
		diaryNoteRouter.POST("/", handlers.CreateDiaryNote)
		diaryNoteRouter.GET("/all", handlers.GetAllDiaryNotes)
		diaryNoteRouter.GET("/some", handlers.GetSomeDiaryNote)
		diaryNoteRouter.GET("/:id", handlers.GetDiaryNote)
		diaryNoteRouter.PUT("/:id", handlers.UpdateDiaryNote)
		diaryNoteRouter.DELETE("/:id", handlers.DeleteDiaryNote)
	}
}
