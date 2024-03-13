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

	verticalContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	selectTab := widget.NewSelect(nil, nil)
	for tab := range filter.Group {
		selectTab.Options = append(selectTab.Options, tab)
	}
	if len(selectTab.Options) > 0 {
		selectTab.Selected = selectTab.Options[0]
	}

	selectType := widget.NewSelect(nil, nil)

	updateTable := func(disciplines []string) {
		newTable := createTable(disciplines)

		verticalContainer.Add(newTable)
		verticalContainer.Refresh()
	}

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

		for _, dis := range filter.Group[tab] {
			updateTable(dis[selectType.Selected])
		}
	}

	selectType.OnChanged = func(t string) {
		for _, dis := range filter.Group[selectTab.Selected] {
			updateTable(dis[t])
		}
	}

	verticalContainer.Add(selectTab)
	verticalContainer.Add(selectType)
	verticalContainer.Add(widget.NewLabel("ХУЙ"))

	verticalContainer.Refresh()

	return verticalContainer
}

func createTable(disciplines []string) *widget.Table {
	numRows := len(disciplines) + 1
	numCols := len(days) + 1

	columns := make([]string, numCols)
	columns[0] = "Предметы"
	copy(columns[1:], days)

	table := widget.NewTable(
		func() (int, int) {
			return numRows, numCols
		},
		func() fyne.CanvasObject {
			return widget.NewEntry()
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			entry, ok := o.(*widget.Entry)
			if !ok {
				return
			}

			if i.Row == 0 {
				entry.Text = columns[i.Col]
			} else {
				entry.Text = ""
			}
		},
	)

	for i := 0; i < numCols; i++ {
		table.SetColumnWidth(i, 100)
	}
	for i := 0; i < numRows; i++ {
		table.SetRowHeight(i, 30)
	}

	table.CreateRenderer()

	return table
}
