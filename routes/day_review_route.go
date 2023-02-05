package routes

import (
	"remood/handlers"

	"github.com/gin-gonic/gin"
)

func DayReviewRouter(r *gin.RouterGroup) {
	dayReviewRouter := r.Group("day-reviews")
	{
		dayReviewRouter.GET("/day", handlers.GetDayReview)
		dayReviewRouter.GET("/between", handlers.GetDayReviewsBetween)
		dayReviewRouter.GET("/month", handlers.GetDayReviewsInMonth)
	}
}