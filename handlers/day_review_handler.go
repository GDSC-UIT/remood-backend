package handlers

import (
	"net/http"
	"strconv"

	"remood/models"
	"remood/pkg/auth"
	"remood/pkg/utils"

	"github.com/gin-gonic/gin"
)


func GetDayReview(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	day := ctx.Query("day")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Invalid parameter"))
		return
	}

	var dayReview models.DayReview
	dayReview.UserID = claims.ID

	err = dayReview.GetOne(day)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get day review aggregation"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get day review aggregation successfully", 
			gin.H{"day_review": dayReview}))
}

func GetDayReviewsBetween(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	startDayString := ctx.Query("start-day")
	endDayString := ctx.Query("end-day")

	var startDay, endDay int64
	startDay, err = utils.StringToInt64(startDayString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Invalid parameter"))
		return
	}
	endDay, err = utils.StringToInt64(endDayString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Invalid parameter"))
		return
	}


	var dayReview models.DayReview
	dayReview.UserID = claims.ID
	dayReviews, err := dayReview.GetBetween(startDay, endDay)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get day review between two days"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get day reviews between two days successfully",
			gin.H{"day_reviews": dayReviews}))
}

func GetDayReviewsInMonth(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	monthString := ctx.Query("month")

	monthInt, err := strconv.Atoi(monthString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Invalid parameter"))
		return
	}

	month := int64(monthInt)

	var dayReview models.DayReview
	dayReview.UserID = claims.ID
	dayReviews, err := dayReview.GetInMonth(month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get day reviews in month"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get day reviews in month successfully", 
			gin.H{"day_reviews": dayReviews}))
}