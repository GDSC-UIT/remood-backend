package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func DayReviewRouter(r *gin.RouterGroup) {
	dayReviewRouter := r.Group("day-reviews")
	{
		dayReviewRouter.GET("/all", handlers.GetAllDayReviews)
		dayReviewRouter.GET("/some", handlers.GetSomeDayReview)
		dayReviewRouter.GET("/:id", handlers.GetDayReview)
	}
}