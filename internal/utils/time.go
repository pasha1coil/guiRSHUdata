package utils

import "time"

func GetCurrentMonth() string {
	return time.Now().Month().String()
}

func GetPreviousMonth() string {
	return time.Now().AddDate(0, -1, 0).Month().String()
}
