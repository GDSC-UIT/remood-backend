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

func CreateDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var diaryNote models.ReviewNote
	diaryNote.UserID = claims.ID

	if err := ctx.BindJSON(&diaryNote); err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read diary note info"))
		return
	}

	if err := diaryNote.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to create diary note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Create diary note successfully", 
			gin.H{"diary_note": diaryNote}))
}


func CreateManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var diaryNotes []models.DiaryNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read diary notes info"))
		return
	}

	err = json.Unmarshal(body, &diaryNotes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read diary notes info"))
		return
	}

	var d models.DiaryNote
	d.UserID = claims.ID
	if err = d.CreateMany(diaryNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to create diary notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Create diary notes successfully", gin.H{
			"diary_notes": diaryNotes,
		},))
}

func GetAllDiaryNotes(ctx *gin.Context) {
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

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID
	diaryNotes, err := diaryNote.GetAll(sort_by_time, filter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get diary notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get all diary notes successfully",
			gin.H{"diary_notes": diaryNotes}))
}

func GetDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	ID := ctx.Query("id")

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	if err = diaryNote.GetOne(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get diary note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get diary note successfully", 
			gin.H{"diary_note": diaryNote}))
}

func GetSomeDiaryNotes(ctx *gin.Context) {
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

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID
	diaryNotes, err := diaryNote.GetSome(page, limit, sort_by_time, filter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get diary notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get some diary notes successfully", 
			gin.H{"diary_notes": diaryNotes}))
}

func UpdateDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var newDiaryNote models.DiaryNote
	if err := ctx.BindJSON(&newDiaryNote); err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read new diary note info"))
		return
	}

	if claims.ID != newDiaryNote.UserID {
		ctx.JSON(http.StatusForbidden, 
			models.ErrorResponse("Invalid User ID"))
		return
	}

	var diaryNote models.DiaryNote
	if err := diaryNote.Update(newDiaryNote); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to update diary note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Update diary note successfully", nil))
}

func UpdateManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var diaryNotes []models.DiaryNote
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read diary notes info"))
		return
	}

	err = json.Unmarshal(body, &diaryNotes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read diary notes info"))
		return
	}

	// Check valid UserID of each diary note
	for _, d := range diaryNotes {
		if d.UserID != claims.ID {
			ctx.JSON(http.StatusForbidden, 
				models.ErrorResponse("Invalid UserID of diary notes"))
			return
		}
	}

	var d models.DiaryNote
	if err = d.UpdateMany(diaryNotes); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to some diary notes info"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Update diary notes successfully", nil))
}

func PinDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := diaryNote.Pin(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to pin diary note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Pin diary note successfully", nil))
}

func DeleteDiaryNote(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var diaryNote models.DiaryNote
	diaryNote.UserID = claims.ID

	ID := ctx.Query("id")

	if err := diaryNote.Delete(ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to delete diary note"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Delete diary note successfully", nil))
}

func DeleteManyDiaryNotes(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	params := ctx.Query("ids")
	IDs := strings.Split(params, ",")

	var d models.DiaryNote
	d.UserID = claims.ID
	if err = d.DeleteMany(IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to delete some diary notes"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Delete diary notes succesfully", nil))
}