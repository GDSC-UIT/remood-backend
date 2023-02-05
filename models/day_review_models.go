package models

import (
	"context"
	"strconv"
	// "log"
	"time"

	"remood/pkg/const/collections"
	"remood/pkg/database"
	"remood/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Point of day review is aggregation of all review notes of that day
type DayReview struct {
	UserID primitive.ObjectID `json:"user_id"`
	Date   int64              `json:"date"`
	Point  float32            `json:"point"`
}

// Get day review of a day
func (dr *DayReview) GetOne(day string) error {
	var r ReviewNote
	r.UserID = dr.UserID

	reviewNotes, err := r.GetAll("", gin.H{"day": day})
	if err != nil {
		return err
	}

	point := dr.AggregateReviewNoteFromSlice(reviewNotes, 0, len(reviewNotes)-1)
	dr.Point = point
	dayInt, _ := strconv.Atoi(day)
	dr.Date = int64(dayInt)

	return nil
}


// Get day reviews between two days
func (dr *DayReview) GetBetween(start int64, end int64) ([]DayReview, error) {
	filter, err := utils.MakeBetweenFilter(start, end)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	utils.SetSortForFindOption(opts, "asc")
	collections := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))
	cursor, err := collections.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var reviewNotes []ReviewNote
	if err = cursor.All(context.Background(), &reviewNotes); err != nil {
		return nil, err
	}

	// Aggregate review notes by day
	result := make([]DayReview, 0)
	startDay := utils.GetDayFromInt64(start)
	endDay := utils.GetDayFromInt64(end)
	var i int = 0
	var j int
	var point float32

	numberOfReviewNotes := len(reviewNotes)
	for utils.IsSmallerOrEqualDate(startDay, endDay) {
		j = dr.GetLastReviewNoteOfDayIndex(reviewNotes, startDay, i, numberOfReviewNotes-1)
		if j >= i {
			point = dr.AggregateReviewNoteFromSlice(reviewNotes, i, j)
			i = j + 1
		} else {
			point = -1
		}

		result = append(result, DayReview{
			UserID: dr.UserID,
			Date:   startDay.Unix(),
			Point:  point,
		})
		startDay = startDay.AddDate(0, 0, 1)
	}
	return result, nil
}

// Get day review of days in a month
func (dr *DayReview) GetInMonth(month int64) ([]DayReview, error) {
	startDay, endDay := utils.GetStartAndEndDayOfMonth(month)
	return dr.GetBetween(startDay, endDay)
}





// output: the last index of review note in expected day
// if it < startIndex, no review note in that day
func (dr *DayReview) GetLastReviewNoteOfDayIndex(reviewNotes []ReviewNote, day time.Time, startIndex int, endIndex int) int {
	for ;startIndex <= endIndex; startIndex++ {
		reviewDay := utils.GetDayFromInt64(reviewNotes[startIndex].CreatedAt)
		if reviewDay != day {
			break
		}
	}
	return startIndex - 1
}

// Calculage aggregation by medium of all point of a day
func (dr *DayReview) AggregateReviewNoteFromSlice(reviewNotes []ReviewNote, start int, end int) float32 {
	result := float32(0)
	if end < start {
		return 0
	}
	for i := start; i <= end; i++ {
		result += float32(reviewNotes[i].Point)
	}
	result /= float32(end - start + 1)
	return result
}