package service

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/tealeg/xlsx"
	"sort"
	"strconv"
)

func (s *Service) generateReport(finishData map[string][]models.EntryData) {
	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Отчет")
	if err != nil {
		errorMessage := "Ошибка при создании листа таблицы:" + err.Error()
		dialog.ShowError(errors.New(errorMessage), models.TopWindow)
	}

	headers := []string{"№ п/п", "Название предмета", "Дата", "День недели", "Тип недели", "Факультет, курс, группа", "Тип занятий", "Часы"}
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}

	upCh := make(chan map[string]models.GenerateReport)
	downCh := make(chan map[string]models.GenerateReport)

	go func() {
		upMap := make(map[string]models.GenerateReport)
		for _, entries := range finishData {
			for _, entry := range entries {
				for _, dayMap := range entry.UpperDay {
					for _, subjects := range dayMap {
						for _, subject := range subjects {
							hash := utils.GenerateHash(entry.Month + entry.Type + entry.Group + subject.Number + subject.Subject + subject.WeekDay)
							if data, ok := upMap[hash]; ok {
								if data.Created.Before(subject.Created) {
									upMap[hash] = models.GenerateReport{
										Month:    entry.Month,
										Group:    entry.Group,
										Type:     entry.Type,
										Subject:  subject.Subject,
										DayWeek:  subject.WeekDay,
										Number:   subject.Number,
										Entry:    subject.Entry,
										Created:  subject.Created,
										TypeWeek: models.UpperWeek,
									}
								}
							} else {
								upMap[hash] = models.GenerateReport{
									Month:    entry.Month,
									Group:    entry.Group,
									Type:     entry.Type,
									Subject:  subject.Subject,
									DayWeek:  subject.WeekDay,
									Number:   subject.Number,
									Entry:    subject.Entry,
									Created:  subject.Created,
									TypeWeek: models.UpperWeek,
								}
							}
						}
					}
				}
			}
		}
		upCh <- upMap
	}()

	go func() {
		downMap := make(map[string]models.GenerateReport)
		for _, entries := range finishData {
			for _, entry := range entries {
				for _, dayMap := range entry.LowerDay {
					for _, subjects := range dayMap {
						for _, subject := range subjects {
							hash := utils.GenerateHash(entry.Month + entry.Type + entry.Group + subject.Number + subject.Subject + subject.WeekDay)
							if data, ok := downMap[hash]; ok {
								if data.Created.Before(subject.Created) {
									downMap[hash] = models.GenerateReport{
										Month:    entry.Month,
										Group:    entry.Group,
										Type:     entry.Type,
										Subject:  subject.Subject,
										DayWeek:  subject.WeekDay,
										Number:   subject.Number,
										Entry:    subject.Entry,
										Created:  subject.Created,
										TypeWeek: models.LowerWeek,
									}
								}
							} else {
								downMap[hash] = models.GenerateReport{
									Month:    entry.Month,
									Group:    entry.Group,
									Type:     entry.Type,
									Subject:  subject.Subject,
									DayWeek:  subject.WeekDay,
									Number:   subject.Number,
									Entry:    subject.Entry,
									Created:  subject.Created,
									TypeWeek: models.LowerWeek,
								}
							}
						}
					}
				}
			}
		}
		downCh <- downMap
	}()

	forGenerateUp := <-upCh
	forGenerateDown := <-downCh

	sort.Slice(models.DaysInfo, func(i, j int) bool {
		return models.DaysInfo[i].Date.Before(models.DaysInfo[j].Date)
	})

	hoursByNumber := make(map[string]float64)

	for _, day := range models.DaysInfo {
		for _, res := range forGenerateUp {
			if day.Date.Month().String() == res.Month && day.Weekday.String() == res.DayWeek && day.WeekType == res.TypeWeek {
				row := sheet.AddRow()

				cell := row.AddCell()
				cell.Value = res.Number

				cell = row.AddCell()
				cell.Value = res.Subject

				cell = row.AddCell()
				cell.Value = day.Date.Format("2006-01-02")

				cell = row.AddCell()
				cell.Value = models.RussianWeekday[res.DayWeek]

				cell = row.AddCell()
				cell.Value = "Верхняя"

				cell = row.AddCell()
				cell.Value = res.Group

				cell = row.AddCell()
				cell.Value = res.Type

				cell = row.AddCell()
				cell.Value = res.Entry

				hours := hoursByNumber[res.Number]
				intEntry, err := strconv.ParseFloat(res.Entry, 2)
				if err != nil {
					intEntry = 0
				}
				hoursByNumber[res.Number] = hours + intEntry

			}
		}

		for _, res := range forGenerateDown {
			if day.Date.Month().String() == res.Month && day.Weekday.String() == res.DayWeek && day.WeekType == res.TypeWeek {
				row := sheet.AddRow()

				cell := row.AddCell()
				cell.Value = res.Number

				cell = row.AddCell()
				cell.Value = res.Subject

				cell = row.AddCell()
				cell.Value = day.Date.Format("2006-01-02")

				cell = row.AddCell()
				cell.Value = models.RussianWeekday[res.DayWeek]

				cell = row.AddCell()
				cell.Value = "Нижняя"

				cell = row.AddCell()
				cell.Value = res.Group

				cell = row.AddCell()
				cell.Value = res.Type

				cell = row.AddCell()
				cell.Value = res.Entry

				hours := hoursByNumber[res.Number]
				intEntry, err := strconv.ParseFloat(res.Entry, 2)
				if err != nil {
					intEntry = 0
				}
				hoursByNumber[res.Number] = hours + intEntry
			}
		}
	}

	for i := 0; i < 5; i++ {
		row := sheet.AddRow()
		if i == 4 {
			cell := row.AddCell()
			cell.Value = "№ п/п"
			cell = row.AddCell()
			cell.Value = "Общее количество часов"
		}
	}

	for number, hours := range hoursByNumber {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = number
		cell = row.AddCell()
		cell.SetFloat(hours)
	}

	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			dialog.ShowError(errors.New(fmt.Sprintf("ошибка при выборе папки: %s", err)), models.TopWindow)
			return
		}

		folderPath := uri.Path()

		reportFilePath := folderPath + "/report.xlsx"
		if err := file.Save(reportFilePath); err != nil {
			dialog.ShowError(errors.New(fmt.Sprintf("ошибка при сохранении файла: %s", err)), models.TopWindow)
			return
		}

		dialog.ShowInformation("Сохранено", "Файл отчета успешно сохранен в "+reportFilePath, models.TopWindow)
	}, models.TopWindow)

	FinishData = make(map[string][]models.EntryData)
}
