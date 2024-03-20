package models

import "time"

var Month = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

type WeekType int

const (
	UpperWeek WeekType = iota
	LowerWeek
)

type DayInfo struct {
	Date     time.Time
	Weekday  time.Weekday
	WeekType WeekType
}
