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
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	if err := ctx.BindJSON(&diaryNote); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read diary note info",
			Error:   true,
		})
		return
	}

	if err := diaryNote.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to create diary note",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Create diary note successfully",
		Error:   false,
		Data: gin.H{
			"diary_note": diaryNote,
		},
	})
}

func CreateManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var diaryNotes []models.DiaryNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read diary notes info",
			Error:   true,
		})
		return
	}

	err = json.Unmarshal(body, &diaryNotes)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read diary notes info",
			Error:   true,
		})
		return
	}

	var d models.DiaryNote
	d.UserID = claims.ID
	if err = d.CreateMany(diaryNotes); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to create diary notes",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Create diary notes successfully",
		Error:   false,
		Data: gin.H{
			"diary_notes": diaryNotes,
		},
	})
}

func GetAllDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
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
			ctx.JSON(http.StatusBadRequest, models.Response{
				Message: "Fail to read parameters",
				Error:   true,
			})
			return
		}

		filter = make(gin.H)
		for i := range filter_bys {
			filter[filter_bys[i]] = filter_values[i]
		}
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID
	diaryNotes, err := diaryNote.GetAll(sort_by_time, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to get diary notes",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Get all diary notes successfully",
		Error:   false,
		Data: gin.H{
			"diary_notes": diaryNotes,
		},
	})
}

func GetDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	ID := ctx.Query("id")

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	if err = diaryNote.GetOne(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to get diary note",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Get diary note successfully",
		Error:   false,
		Data: gin.H{
			"diary_note": diaryNote,
		},
	})
}

func GetSomeDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
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
			ctx.JSON(http.StatusBadRequest, models.Response{
				Message: "Fail to read parameters",
				Error:   true,
			})
			return
		}

		filter = make(gin.H)
		for i := range filter_bys {
			filter[filter_bys[i]] = filter_values[i]
		}
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID
	diaryNotes, err := diaryNote.GetSome(page, limit, sort_by_time, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to get diary notes",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Get some diary notes successfully",
		Error:   false,
		Data: gin.H{
			"diary_notes": diaryNotes,
		},
	})
}

func UpdateDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var newDiaryNote models.DiaryNote
	if err := ctx.BindJSON(&newDiaryNote); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read new diary note info",
			Error:   true,
		})
		return
	}

	if claims.ID != newDiaryNote.UserID {
		ctx.JSON(http.StatusForbidden, models.Response{
			Message: "Invalid User ID",
			Error:   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	if err := diaryNote.Update(newDiaryNote); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to update diary note",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Update diary note successfully",
		Error:   false,
	})
}

func UpdateManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var diaryNotes []models.DiaryNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read diary notes info",
			Error:   true,
		})
		return
	}

	err = json.Unmarshal(body, &diaryNotes)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read diary notes info",
			Error:   true,
		})
		return
	}

	// Check valid UserID of each diary note
	for _, d := range diaryNotes {
		if d.UserID != claims.ID {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Message: "Invalid UserID of some diary notes",
				Error:   true,
			})
			return
		}
	}

	var d models.DiaryNote
	if err = d.UpdateMany(diaryNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to some diary notes info",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Update diary notes successfully",
		Error:   false,
	})
}

func PinDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := diaryNote.Pin(ID); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to pin diary note",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Pin diary note successfully",
		Error:   false,
	})
}

func DeleteDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := diaryNote.Delete(ID); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to delete diary note",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Delete diary note successfully",
		Error:   false,
	})
}

func DeleteManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	params := ctx.Query("ids")
	IDs := strings.Split(params, ",")

	var d models.DiaryNote
	d.UserID = claims.ID
	if err = d.DeleteMany(IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to delete some diary notes",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Delete diary notes succesfully",
		Error:   false,
	})
}


