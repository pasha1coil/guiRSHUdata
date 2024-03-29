package utils

import (
	"demofine/internal/adapters"
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

	productCalendar := adapters.ProductClient(start, end)

	if current.Weekday() != time.Monday && current.Weekday() != time.Sunday {
		for current.Weekday() != time.Monday {
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

	for _, date := range productCalendar.Days {
		for index, ourDate := range days {
			if date.Date == ourDate.Date.Format("02.01.2006") {
				if _, ok := models.SpecialTypes[date.TypeText]; ok {
					days = append(days[:index], days[index+1:]...)
					continue
				}
			}
		}
	}
	return days
}
