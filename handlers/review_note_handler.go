package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID

	if err := ctx.BindJSON(&reviewNote); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read review note info",
			"error":   true,
		})
		return
	}

	if err := reviewNote.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to create review note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create review note successfully",
		"error":   false,
		"data": gin.H{
			"review_note": reviewNote,
		},
	})
}

func CreateManyReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var reviewNotes []models.ReviewNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read review notes info",
			"error":   true,
		})
		return
	}

	err = json.Unmarshal(body, &reviewNotes)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read review notes info",
			"error":   true,
		})
		return
	}

	var d models.ReviewNote
	d.UserID = claims.ID
	if err = d.CreateMany(reviewNotes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to create review notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create review notes successfully",
		"error":   false,
		"data": gin.H{
			"review_notes": reviewNotes,
		},
	})
}

func GetAllReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Fail to read parameters",
				"error":   true,
			})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to get review notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all review notes successfully",
		"error":   false,
		"data": gin.H{
			"review_note": reviewNotes,
		},
	})
}

func GetReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	ID := ctx.Query("id")

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID

	if err = reviewNote.GetOne(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to get review note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get review note successfully",
		"error":   false,
		"data": gin.H{
			"review_notes": reviewNote,
		},
	})
}

func GetSomeReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Fail to read parameters",
				"error":   true,
			})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to get review notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get some review notes successfully",
		"error":   false,
		"data": gin.H{
			"review_notes": reviewNotes,
		},
	})
}

func UpdateReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var newReviewNote models.ReviewNote
	if err := ctx.BindJSON(&newReviewNote); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read new review note info",
			"error":   true,
		})
		return
	}

	if claims.ID != newReviewNote.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid User ID",
			"error":   true,
		})
		return
	}

	var reviewNote models.ReviewNote
	if err := reviewNote.Update(newReviewNote); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to update review note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update review note successfully",
		"error":   false,
	})
}

func UpdateManyReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var reviewNotes []models.ReviewNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read review notes info",
			"error":   true,
		})
		return
	}

	err = json.Unmarshal(body, &reviewNotes)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read review notes info",
			"error":   true,
		})
		return
	}

	// Check valid UserID of each review note
	for _, d := range reviewNotes {
		if d.UserID != claims.ID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid UserID of some review notes",
				"error":   true,
			})
			return
		}
	}

	var d models.ReviewNote
	if err = d.UpdateMany(reviewNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to some review notes info",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update review notes successfully",
		"error":   false,
	})
}

func DeleteReviewNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var reviewNote models.ReviewNote
	reviewNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := reviewNote.Delete(ID); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to delete review note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete review note successfully",
		"error":   false,
	})
}

func DeleteManyReviewNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	params := ctx.Query("ids")
	IDs := strings.Split(params, ",")

	var d models.ReviewNote
	d.UserID = claims.ID
	if err = d.DeleteMany(IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to delete some review notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete review notes succesfully",
		"error":   false,
	})
}
