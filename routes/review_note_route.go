package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func ReviewNoteRouter(r *gin.RouterGroup) {
	reviewNoteRouter := r.Group("review-notes")
	{
		reviewNoteRouter.POST("/", handlers.CreateReviewNote)
		reviewNoteRouter.POST("/many", handlers.CreateManyReviewNotes)
		reviewNoteRouter.GET("/all", handlers.GetAllReviewNotes)
		reviewNoteRouter.GET("/some", handlers.GetSomeReviewNotes)
		reviewNoteRouter.GET("/", handlers.GetReviewNote)
		reviewNoteRouter.PUT("/", handlers.UpdateReviewNote)
		reviewNoteRouter.PUT("/many", handlers.UpdateManyReviewNotes)
		reviewNoteRouter.DELETE("/", handlers.DeleteReviewNote)
		reviewNoteRouter.DELETE("/many", handlers.DeleteManyReviewNotes)
	}
}
