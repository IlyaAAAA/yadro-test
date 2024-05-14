package util

import (
	"strconv"
	"time"
)

func GetTimeStr(time time.Time) string {
	var hours string
	if time.Hour() < 10 {
		hours = "0" + strconv.Itoa(time.Hour())
	} else {
		hours = strconv.Itoa(time.Hour())
	}

	var minutes string
	if time.Minute() < 10 {
		minutes = "0" + strconv.Itoa(time.Minute())
	} else {
		minutes = strconv.Itoa(time.Minute())
	}

	return hours + ":" + minutes
}
