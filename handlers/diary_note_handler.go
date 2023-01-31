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

func CreateDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	if err := ctx.BindJSON(&diaryNote); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read diary note info",
			"error":   true,
		})
		return
	}

	if err := diaryNote.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to create diary note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create diary note successfully",
		"error":   false,
		"data": gin.H{
			"diary_note": diaryNote,
		},
	})
}

func CreateManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var diaryNotes []models.DiaryNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read diary notes info",
			"error":   true,
		})
		return
	}

	err = json.Unmarshal(body, &diaryNotes)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read diary notes info",
			"error":   true,
		})
		return
	}

	var d models.DiaryNote
	d.UserID = claims.ID
	if err = d.CreateMany(diaryNotes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to create diary notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create diary notes successfully",
		"error":   false,
		"data": gin.H{
			"diary_notes": diaryNotes,
		},
	})
}

func GetAllDiaryNotes(ctx *gin.Context) {
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

	var filter map[string]interface{}
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

		filter = make(map[string]interface{})
		for i := range filter_bys {
			filter[filter_bys[i]] = filter_values[i]
		}
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID
	diaryNotes, err := diaryNote.GetAll(sort_by_time, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to get diary notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all diary notes successfully",
		"error":   false,
		"data": gin.H{
			"diary_note": diaryNotes,
		},
	})
}

func GetDiaryNote(ctx *gin.Context) {
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

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	if err = diaryNote.GetOne(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to get diary note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get diary note successfully",
		"error":   false,
		"data": gin.H{
			"diary_notes": diaryNote,
		},
	})
}

func GetSomeDiaryNotes(ctx *gin.Context) {
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

	var filter map[string]interface{}
	// In case no filter
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

		filter = make(map[string]interface{})
		for i := range filter_bys {
			filter[filter_bys[i]] = filter_values[i]
		}
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID
	diaryNotes, err := diaryNote.GetMany(page, limit, sort_by_time, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to get diary notes",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get some diary notes successfully",
		"error":   false,
		"data": gin.H{
			"diary_notes": diaryNotes,
		},
	})
}

func UpdateDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	var newDiaryNote models.DiaryNote
	if err := ctx.BindJSON(&newDiaryNote); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read new diary note info",
			"error":   true,
		})
		return
	}

	if diaryNote.UserID != newDiaryNote.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Can not update other user's diary note",
			"error":   true,
		})
		return
	}

	if err := diaryNote.Update(newDiaryNote); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to update diary note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update diary note successfully",
		"error":   false,
	})
}

func UpdateManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var diaryNotes []models.DiaryNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read diary notes info",
			"error":   true,
		})
		return
	}

	err = json.Unmarshal(body, &diaryNotes)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to read diary notes info",
			"error":   true,
		})
		return
	}

	// Check valid UserID of each diary note
	for _, d := range diaryNotes {
		if d.UserID != claims.ID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid UserID of some diary note",
				"error":   true,
			})
			return
		}
	}

	var d models.DiaryNote
	if err = d.UpdateMany(diaryNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to some diary notes info",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update diary notes successfully",
		"error":   false,
	})
}

func DeleteDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := diaryNote.Delete(ID); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to delete diary note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete diary note successfully",
		"error":   false,
	})
}

func DeleteManyDiaryNote(ctx *gin.Context) {
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

	var d models.DiaryNote
	d.UserID = claims.ID
	if err = d.DeleteMany(IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to delete some diary note",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete diary notes succesfully",
		"error":   false,
	})
}
