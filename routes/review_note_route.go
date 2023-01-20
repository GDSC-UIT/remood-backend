package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func ReviewNoteRouter(r *gin.RouterGroup) {
	reviewNoteRouter := r.Group("review-notes")
	{
		reviewNoteRouter.POST("/", handlers.CreateReviewNote)
		reviewNoteRouter.GET("/all", handlers.GetAllReviewNotes)
		reviewNoteRouter.GET("/some", handlers.GetSomeReviewNote)
		reviewNoteRouter.GET("/:id", handlers.GetReviewNote)
		reviewNoteRouter.PUT("/:id", handlers.UpdateReviewNote)
		reviewNoteRouter.DELETE("/:id", handlers.DeleteReviewNote)
	}
}
