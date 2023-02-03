package utils

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetSortForFindOption(opts *options.FindOptions, sort_by_time string) {
	if sort_by_time == "asc" {
		opts.SetSort(bson.M{"created_at": 1})
	} else if sort_by_time == "desc" {
		opts.SetSort(bson.M{"created_at": -1})
	}
}

func SetSkipPage(opts *options.FindOptions, page int64, limit int64) {
	opts.SetSkip(page * limit).SetLimit(limit)
}

func MakeFilter(rawFilter gin.H) (bson.M, error) {
	var filter bson.M
	if rawFilter != nil {
		filter = bson.M(rawFilter)

		// Day/Month filter
		if rawFilter["day"] != nil {
			// Convert filter time to date time format
			filterTime, err := StringToInt64(rawFilter["day"].(string))
			if err != nil {
				return nil, err
			}
			startOfDay := GetDayFromInt64(filterTime)
			endOfDay := startOfDay.AddDate(0, 0, 1)

			filter["created_at"] = bson.M{
				"$gte": startOfDay.Unix(),
				"$lte": endOfDay.Unix(),
			}
			delete(filter, "day")
		} else if rawFilter["month"] != nil {
			// Convert filter time to date time format
			filterTime, err := StringToInt64(rawFilter["month"].(string))
			if err != nil {
				return nil, err
			}
			startOfMonth := GetDayFromInt64(filterTime)
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			filter["created_at"] = bson.M{
				"$gte": startOfMonth.Unix(),
				"$lte": endOfMonth.Unix(),
			}
			delete(filter, "month")
		}
	}

	return filter, nil
}

func MakeBetweenFilter(startDay int64, endDay int64) (bson.M, error) {
	filter := make(bson.M)
	startTime := GetDayFromInt64(startDay)
	endTime := GetDayFromInt64(endDay)

	filter["created_at"] = bson.M{
		"$gte": startTime.Unix(),
		"$lte": endTime.AddDate(0, 0, 1).Unix(),
	}

	return filter, nil
}