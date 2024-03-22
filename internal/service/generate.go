package service

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/tealeg/xlsx"
	"log"
)

func (s *Service) generateReport(finishData map[string][]models.EntryData) {
	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Отчет")
	if err != nil {
		log.Print("Ошибка при создании листа:", err)
		return
	}

	headers := []string{"№ п/п", "Название предмета", "Дата", "День недели", "Тип недели", "Факультет, курс, группа", "Тип занятий", "Часы"}
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}

	forGenerateUp := make(map[string]models.GenerateReport)
	forGenerateDown := make(map[string]models.GenerateReport)

	for _, entries := range finishData {
		for _, entry := range entries {
			for _, dayMap := range entry.UpperDay {
				for _, subjects := range dayMap {
					for _, subject := range subjects {
						hash := utils.GenerateHash(entry.Month + entry.Type + entry.Group + subject.Number + subject.Subject + subject.WeekDay)
						if data, ok := forGenerateUp[hash]; ok {
							if data.Created.Before(subject.Created) {
								forGenerateUp[hash] = models.GenerateReport{
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
							forGenerateUp[hash] = models.GenerateReport{
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
			for _, dayMap := range entry.LowerDay {
				for _, subjects := range dayMap {
					for _, subject := range subjects {
						hash := utils.GenerateHash(entry.Month + entry.Type + entry.Group + subject.Number + subject.Subject + subject.WeekDay)
						if data, ok := forGenerateDown[hash]; ok {
							if data.Created.Before(subject.Created) {
								forGenerateDown[hash] = models.GenerateReport{
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
							forGenerateDown[hash] = models.GenerateReport{
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
			}
		}
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
}
