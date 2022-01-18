package utils

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

func GetTimeArr(start, end string) int64 {
	timeLayout := "2016/01/02"
	loc, _ := time.LoadLocation("Local")
	// 转成时间戳
	startUnix, _ := time.ParseInLocation(timeLayout, start, loc)
	endUnix, _ := time.ParseInLocation(timeLayout, end, loc)
	startTime := startUnix.Unix()
	endTime := endUnix.Unix()
	// 求相差天数
	date := (endTime - startTime) / 86400
	return date
}

func ParseTimestampToTime(ts *timestamp.Timestamp) *time.Time {
	if ts == nil || ts.Seconds == 0 {
		return nil
	}
	tt, _ := ptypes.Timestamp(ts)
	return &tt
}
func ParseTimeToTimestamp(t *time.Time) *timestamp.Timestamp {
	if t == nil {
		return nil
	}
	ts, _ := ptypes.TimestampProto(*t)
	return ts
}
