package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func DiaryNoteRouter(r *gin.RouterGroup) {
	diaryNoteRouter := r.Group("diary-notes")
	{
		diaryNoteRouter.POST("/", handlers.CreateDiaryNote)
		diaryNoteRouter.POST("/many", handlers.CreateManyDiaryNotes)
		diaryNoteRouter.GET("/all", handlers.GetAllDiaryNotes)
		diaryNoteRouter.GET("/some", handlers.GetSomeDiaryNotes)
		diaryNoteRouter.GET("/", handlers.GetDiaryNote)
		diaryNoteRouter.PUT("/", handlers.UpdateDiaryNote)
		diaryNoteRouter.PUT("/many", handlers.UpdateManyDiaryNotes)
		diaryNoteRouter.DELETE("/", handlers.DeleteDiaryNote)
		diaryNoteRouter.DELETE("/many", handlers.DeleteManyDiaryNote)
	}
}
