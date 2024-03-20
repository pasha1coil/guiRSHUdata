package service

import (
	"demofine/internal/models"
	"time"
)

func GenerateDays() []models.DayInfo {
	days := []models.DayInfo{}
	var start, end time.Time
	cu := time.Now()
	mo := models.Month[cu.Month().String()]

	if mo >= 9 {
		start = time.Date(time.Now().Year(), time.September, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(time.Now().Year()+1, time.September, 1, 0, 0, 0, 0, time.UTC)
	} else {
		start = time.Date(time.Now().Year()-1, time.September, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(time.Now().Year(), time.September, 1, 0, 0, 0, 0, time.UTC)
	}

	current := start
	startWeekType := models.UpperWeek

	if current.Weekday() != time.Monday && current.Weekday() != time.Sunday {
		for current.Weekday() != time.Sunday {
			days = append(days, models.DayInfo{
				Date:     current,
				Weekday:  current.Weekday(),
				WeekType: models.UpperWeek,
			})

			current = current.AddDate(0, 0, 1)
		}
		startWeekType = models.LowerWeek
	} else if current.Weekday() == time.Sunday {

		days = append(days, models.DayInfo{
			Date:     current,
			Weekday:  current.Weekday(),
			WeekType: models.UpperWeek,
		})
		current = current.AddDate(0, 0, 1)
		startWeekType = models.LowerWeek
	}

	for current.Before(end) {
		for i := 0; i < 7; i++ {
			weekday := current.Weekday()

			days = append(days, models.DayInfo{
				Date:     current,
				Weekday:  weekday,
				WeekType: startWeekType,
			})
			current = current.AddDate(0, 0, 1)
		}
		if startWeekType == models.UpperWeek {
			startWeekType = models.LowerWeek
		} else {
			startWeekType = models.UpperWeek
		}
	}
	return days
}
