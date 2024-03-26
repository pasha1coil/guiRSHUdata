package models

import "time"

// Token доступа для https://production-calendar.ru
const Token = "3cb4909cf48fd4dbb961a0947605e4f8"
const Country = "ru"

// пример - https://production-calendar.ru/get-period/6914a408120146bcb82ab95c003bc6ad/ru/01092023-31082024/json
const ProductUrl = "https://production-calendar.ru/get-period/"

var SpecialTypes = map[string]struct{}{
	"Выходной день":                               {},
	"Государственный праздник":                    {},
	"Региональный праздник":                       {},
	"Предпраздничный сокращенный рабочий день":    {},
	"Дополнительный / перенесенный выходной день": {},
}

type CustomDate struct {
	time.Time
}

func (c *CustomDate) UnmarshalJSON(b []byte) error {
	dateStr := string(b)
	dateStr = dateStr[1 : len(dateStr)-1] // Убираем кавычки из строки
	t, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return err
	}
	c.Time = t
	return nil
}

type ProductCalendarInfo struct {
	Status       string     `json:"status"`
	CountryCode  string     `json:"country_code"`
	CountryText  string     `json:"country_text"`
	StartDate    CustomDate `json:"dt_start"`
	EndDate      CustomDate `json:"dt_end"`
	WorkWeekType string     `json:"work_week_type"`
	Period       string     `json:"period"`
	Note         string     `json:"note"`
	Days         []DInfo    `json:"days"`
	Statistic    Statistic  `json:"statistic"`
}

type DInfo struct {
	Date         string `json:"date"`
	TypeID       int    `json:"type_id"`
	TypeText     string `json:"type_text"`
	WeekDay      string `json:"week_day"`
	WorkingHours int    `json:"working_hours"`
}

type Statistic struct {
	CalendarDays                int `json:"calendar_days"`
	CalendarDaysWithoutHolidays int `json:"calendar_days_without_holidays"`
	WorkDays                    int `json:"work_days"`
	Weekends                    int `json:"weekends"`
	Holidays                    int `json:"holidays"`
	WorkingHours                int `json:"working_hours"`
}
