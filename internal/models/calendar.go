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

var RussianMonth = map[string]string{
	"January":   "Январь",
	"February":  "Февраль",
	"March":     "Март",
	"April":     "Апрель",
	"May":       "Май",
	"June":      "Июнь",
	"July":      "Июль",
	"August":    "Август",
	"September": "Сентябрь",
	"October":   "Октябрь",
	"November":  "Ноябрь",
	"December":  "Декабрь",
}

var RussianWeekday = map[string]string{
	"Monday":    "Понедельник",
	"Tuesday":   "Вторник",
	"Wednesday": "Среда",
	"Thursday":  "Четверг",
	"Friday":    "Пятница",
	"Saturday":  "Суббота",
	"Sunday":    "Воскресенье",
}

var EnglishWeekday = map[string]string{
	"Понедельник": "Monday",
	"Вторник":     "Tuesday",
	"Среда":       "Wednesday",
	"Четверг":     "Thursday",
	"Пятница":     "Friday",
	"Суббота":     "Saturday",
	"Воскресенье": "Sunday",
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
