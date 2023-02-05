package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"remood/models"
	"remood/pkg/auth"

	"github.com/gin-gonic/gin"
)

func CreateReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID

	if err := ctx.BindJSON(&reviewNote); err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read review note info"))
		return
	}

	if err := reviewNote.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to create review note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Create review note successfully", 
			gin.H{"review_note": reviewNote}))
}

func CreateManyReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var reviewNotes []models.ReviewNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read review notes info"))
		return
	}

	err = json.Unmarshal(body, &reviewNotes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read review notes info"))
		return
	}

	var d models.ReviewNote
	d.UserID = claims.ID
	if err = d.CreateMany(reviewNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to create review notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Create review notes successfully", gin.H{
			"review_notes": reviewNotes,
		},))
}

func GetAllReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	sort_by_time := ctx.Query("sort-by-time")
	filter_by := ctx.Query("filter-by")
	filter_value := ctx.Query("filter-value")

	var filter gin.H
	// In case having filter
	if filter_by != "" {
		filter_bys := strings.Split(filter_by, ",")
		filter_values := strings.Split(filter_value, ",")

		if len(filter_bys) != len(filter_values) {
			ctx.JSON(http.StatusBadRequest, 
				models.ErrorResponse("Fail to read parameters"))
			return
		}

		filter = make(gin.H)
		for i := range filter_bys {
			filter[filter_bys[i]] = filter_values[i]
		}
	}

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID
	reviewNotes, err := reviewNote.GetAll(sort_by_time, filter)
	
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get review notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get all review notes successfully",
			gin.H{"review_notes": reviewNotes}))
}

func GetReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	ID := ctx.Query("id")

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID

	if err = reviewNote.GetOne(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get review note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get review note successfully", 
			gin.H{"review_note": reviewNote}))
}

func GetSomeReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	pageInt32, _ := strconv.Atoi(ctx.Query("page"))
	limitInt32, _ := strconv.Atoi(ctx.Query("limit"))
	sort_by_time := ctx.Query("sort-by-time")
	filter_by := ctx.Query("filter-by")
	filter_value := ctx.Query("filter-value")

	page := int64(pageInt32)
	limit := int64(limitInt32)

	var filter gin.H
	// In case having filter
	if filter_by != "" {
		filter_bys := strings.Split(filter_by, ",")
		filter_values := strings.Split(filter_value, ",")

		if len(filter_bys) != len(filter_values) {
			ctx.JSON(http.StatusBadRequest, 
				models.ErrorResponse("Fail to read parameters"))
			return
		}

		filter = make(gin.H)
		for i := range filter_bys {
			filter[filter_bys[i]] = filter_values[i]
		}
	}

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID
	reviewNotes, err := reviewNote.GetSome(page, limit, sort_by_time, filter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get review notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get some review notes successfully", 
			gin.H{"review_notes": reviewNotes}))
}

func UpdateReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var newReviewNote models.ReviewNote
	if err := ctx.BindJSON(&newReviewNote); err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read new review note info"))
		return
	}

	if claims.ID != newReviewNote.UserID {
		ctx.JSON(http.StatusForbidden, 
			models.ErrorResponse("Invalid User ID"))
		return
	}

	var reviewNote models.ReviewNote
	if err := reviewNote.Update(newReviewNote); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to update review note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Update review note successfully", nil))
}

func UpdateManyReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var reviewNotes []models.ReviewNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read review notes info"))
		return
	}

	err = json.Unmarshal(body, &reviewNotes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read review notes info"))
		return
	}

	// Check valid UserID of each review note
	for _, d := range reviewNotes {
		if d.UserID != claims.ID {
			ctx.JSON(http.StatusForbidden, 
				models.ErrorResponse("Invalid UserID of review notes"))
			return
		}
	}

	var d models.ReviewNote
	if err = d.UpdateMany(reviewNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to some review notes info"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Update review notes successfully", nil))
}

func DeleteReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := reviewNote.Delete(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to delete review note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Delete review note successfully", nil))
}

func DeleteManyReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	params := ctx.Query("ids")
	IDs := strings.Split(params, ",")

	var d models.ReviewNote
	d.UserID = claims.ID
	if err = d.DeleteMany(IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to delete some review notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Delete review notes succesfully", nil))
}


