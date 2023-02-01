package utils 

import (
	"strconv"
	"time"


	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetSortForFindOption(opts *options.FindOptions,sort_by_time string)  {
	if sort_by_time == "asc" {
		opts.SetSort(bson.M{"created_at": 1})
	} else if sort_by_time == "desc" {
		opts.SetSort(bson.M{"created_at": -1})
	}
}


func MakeBsonFilter(filter gin.H) (bson.M, error) {
	var bsonFilter bson.M
	if filter != nil {
		bsonFilter = bson.M(filter)

		// Day/Month filter
		if bsonFilter["day"] != nil {
			// Get filter time in int32
			filterTimeInt32, err := strconv.Atoi(bsonFilter["day"].(string))
			if err != nil {
				return nil, err
			}

			// Convert filter time to date time format
			filterTime := int64(filterTimeInt32)
			formatTime := time.Unix(filterTime, 0)

			startOfDay := time.Date(formatTime.Year(), formatTime.Month(), formatTime.Day(), 0, 0, 0, 0, time.UTC)
			endOfDay := startOfDay.AddDate(0, 0, 1)

			filter["created_at"] = bson.M{
				"$gte": startOfDay.Unix(),
				"$lt":  endOfDay.Unix(),
			}
			delete(bsonFilter, "day")
		} else if bsonFilter["month"] != nil {
			// Get filter time in int32
			filterTimeInt32, err := strconv.Atoi(bsonFilter["month"].(string))
			if err != nil {
				return nil, err
			}

			// Convert filter time to date time format
			filterTime := int64(filterTimeInt32)
			formatTime := time.Unix(filterTime, 0)

			startOfMonth := time.Date(formatTime.Year(), formatTime.Month(), 1, 0, 0, 0, 0, time.UTC)
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			filter["created_at"] = bson.M{
				"$gte": startOfMonth.Unix(),
				"$lt":  endOfMonth.Unix(),
			}
			delete(bsonFilter, "month")
		}
	}

	return bsonFilter, nil
}

func SetSkipPage(opts *options.FindOptions, page int64, limit int64) {
	opts.SetSkip(page * limit).SetLimit(limit)
}