package utils

import (
	"strconv"
	"time"
)

func IsSmallerOrEqualDate(a time.Time, b time.Time) bool {
	return a.Unix() <= b.Unix()
}

func GetStartAndEndDayOfMonth(month int64) (int64, int64) {
	dayInMonth := time.Unix(month, 0)
	startDay := time.Date(dayInMonth.Year(), dayInMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDay := startDay.AddDate(0, 1, 0)
	return startDay.Unix(), endDay.Unix()
}

func StringToInt64(s string) (int64, error) {
	i32, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int64(i32), nil
}

func GetDayFromInt64(aTimeInDay int64) time.Time {
	aTimeInDayformat := time.Unix(aTimeInDay, 0)
	date := time.Date(aTimeInDayformat.Year(), aTimeInDayformat.Month(), aTimeInDayformat.Day(), 0, 0, 0, 0, time.UTC)
	return date
}