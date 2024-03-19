package service

import (
	"demofine/internal/models"
	"github.com/tealeg/xlsx"
	"log"
)

func (s *Service) generateReport(finishData map[string][]models.EntryData) {
	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Отчет")
	if err != nil {
		log.Fatal("Ошибка при создании листа:", err)
	}

	headers := []string{"№ п/п", "Название предмета", "Дата", "Факультет, курс, группа", "Тип занятий", "Часы"}
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}

	for _, entries := range finishData {
		for _, entry := range entries {
			for _, dayMap := range entry.UpperDay {
				for _, hour := range dayMap {
					row := sheet.AddRow()

					cell := row.AddCell()
					cell.Value = hour[0].Number

					cell = row.AddCell()
					cell.Value = hour[0].Subject

					cell = row.AddCell()
					cell.Value = entry.Month

					cell = row.AddCell()
					cell.Value = entry.Group

					cell = row.AddCell()
					cell.Value = entry.Type

					cell = row.AddCell()
					cell.Value = hour[0].Entry
				}
			}
		}
	}

	if err := file.Save("report.xlsx"); err != nil {
		log.Fatal("Ошибка при сохранении файла отчета:", err)
	}
}
