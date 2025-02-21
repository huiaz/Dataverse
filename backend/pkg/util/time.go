package util

import "time"

func FormatTimeToInt64(dateTime time.Time) int64 {
	return dateTime.Unix()
}
