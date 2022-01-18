package utils

import (
	"time"
)

func ParseStrToDateTime(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	const TIME_LAYOUT = "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	// 转成时间戳
	timestamp1, _ := time.ParseInLocation(TIME_LAYOUT,  dateStr,  loc)
	return &timestamp1, nil
}
func ParseStrToDateTime2(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	const TIME_LAYOUT = "2006-01-02T15:04:05Z"
	loc, _ := time.LoadLocation("Local")
	// 转成时间戳
	timestamp1, _ := time.ParseInLocation(TIME_LAYOUT,  dateStr,  loc)
	return &timestamp1, nil
}
