package service

import (
	"demofine/internal/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tealeg/xlsx"
	"log"
)

func (s *Service) MakeTableTab(_ fyne.Window) fyne.CanvasObject {
	fileData, userName, err := s.LoadDataFromBadger()
	if err != nil {
		log.Println("Error loading data from Badger:", err)
		return nil
	}

	var loadedFile models.LoadedFile

	xlFile, err := xlsx.OpenBinary(fileData)
	if err != nil {
		log.Println("Error opening XLSX file:", err)
		return nil
	}

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			cells := row.Cells

			if len(cells) >= 14 {
				loadedFile.Number = append(loadedFile.Number, cells[0].String())
				loadedFile.Plan = append(loadedFile.Plan, cells[1].Value)
				loadedFile.Faculty = append(loadedFile.Faculty, cells[2].String())
				loadedFile.Block = append(loadedFile.Block, cells[3].String())
				loadedFile.Discipline = append(loadedFile.Discipline, cells[4].String())
				loadedFile.Semester = append(loadedFile.Semester, cells[5].Value)
				loadedFile.Group = append(loadedFile.Group, cells[6].String())
				loadedFile.CountStudents = append(loadedFile.CountStudents, cells[7].Value)
				loadedFile.Type = append(loadedFile.Type, cells[8].String())
				loadedFile.Hours = append(loadedFile.Hours, cells[9].Value)
				loadedFile.FIO = append(loadedFile.FIO, cells[10].String())
				loadedFile.TypeWork = append(loadedFile.TypeWork, cells[11].String())
				loadedFile.RankClever = append(loadedFile.RankClever, cells[12].String())
				loadedFile.ProgressClever = append(loadedFile.ProgressClever, cells[13].String())
			}
		}
	}

	var filter models.Filter
	filter.FIO = userName
	filter.Group = make(map[string][]map[string][]string)
	for i, teacher := range loadedFile.FIO {
		if teacher == userName {
			group := loadedFile.Group[i].(string)
			discipline := loadedFile.Discipline[i].(string)
			lessonType := loadedFile.Type[i].(string)

			if _, ok := filter.Group[group]; !ok {
				filter.Group[group] = make([]map[string][]string, 0)
			}

			var found bool
			for j := range filter.Group[group] {
				if filter.Group[group][j][lessonType] != nil {
					filter.Group[group][j][lessonType] = append(filter.Group[group][j][lessonType], discipline)
					found = true
				}
			}
			if !found {
				newLesson := make(map[string][]string)
				newLesson[lessonType] = []string{discipline}
				filter.Group[group] = append(filter.Group[group], newLesson)
			}
		}
	}

	fmt.Println(filter)

	table := createTable([]string{})

	return table
}

func createTable(teacherDisciplines []string) fyne.CanvasObject {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	tableTop := container.NewGridWithRows(len(teacherDisciplines))
	for _, discipline := range teacherDisciplines {
		label := widget.NewLabel(discipline)
		tableTop.Add(label)
	}

	tableLeft := container.NewVBox()
	for _, day := range days {
		label := widget.NewLabel(day)
		tableLeft.Add(label)
	}

	tableContent := container.NewGridWithRows(len(days))
	for range days {
		for range teacherDisciplines {
			entry := widget.NewEntry()
			tableContent.Add(entry)
		}
	}

	table := container.NewBorder(nil, nil, tableLeft, nil, container.NewHSplit(tableTop, tableContent))

	return table
}
