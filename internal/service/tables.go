package service

import (
	"demofine/internal/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/tealeg/xlsx"
	"log"
)

var days = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}

func (s *Service) MakeTableTab(w fyne.Window) fyne.CanvasObject {
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

	selectTab := widget.NewSelect(nil, nil)
	for tab := range filter.Group {
		selectTab.Options = append(selectTab.Options, tab)
	}
	if len(selectTab.Options) > 0 {
		selectTab.Selected = selectTab.Options[0]
	}

	selectType := widget.NewSelect(nil, nil)

	dialog := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		selectTab,
		selectType,
		widget.NewButton("Показать таблицу", func() {
			var data []string
			var ok bool
			for _, dis := range filter.Group[selectTab.Selected] {
				if data, ok = dis[selectType.Selected]; ok {
					break
				}
			}

			table := createTable(data)
			tableWindow := fyne.CurrentApp().NewWindow("Таблица")
			tableWindow.SetContent(table)
			tableWindow.Resize(fyne.NewSize(800, 600))
			tableWindow.Show()
		}),
	)

	selectTab.OnChanged = func(tab string) {
		selectType.Options = nil
		for _, lesson := range filter.Group[tab] {
			for t := range lesson {
				selectType.Options = append(selectType.Options, t)
			}
		}
		if len(selectType.Options) > 0 {
			selectType.Selected = selectType.Options[0]
		}
	}

	selectType.OnChanged = func(t string) {}

	return dialog
}

func createTable(arr []string) *widget.Table {
	var data [][]string
	data = append(data, arr)
	for _, date := range days {
		var row []string
		row = append(row, date)
		data = append(data, row)
	}

	numRows := len(data)
	numCols := len(data[0])

	table := widget.NewTable(
		func() (int, int) {
			return numRows, numCols
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label, ok := o.(*widget.Label)
			if !ok {
				return
			}

			label.SetText(data[i.Row][i.Col])
		},
	)

	for i := 0; i < numCols; i++ {
		table.SetColumnWidth(i, 100)
	}
	for i := 0; i < numRows; i++ {
		table.SetRowHeight(i, 30)
	}

	return table
}
