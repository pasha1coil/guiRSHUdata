package service

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/tealeg/xlsx"
	"log"
	"strings"
)

var days = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}

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

			table := createTable(data, selectTab.Selected, selectType.Selected)
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

var subjectInfo = make(map[*widget.Entry]string)

func createTable(subjects []string, group string, lessonType string) fyne.CanvasObject {
	upperW := container.NewWithoutLayout()

	subjectRowY := float32(50)

	downW := container.NewWithoutLayout()
	downW.Move(fyne.NewPos(100, 100))

	was := false

	var entryWidgets []*widget.Entry

	for _, subject := range subjects {
		upperSubjectRow := container.NewWithoutLayout()
		downSubjectRow := container.NewWithoutLayout()

		if !was {
			upLb := widget.NewLabel("Верхняя неделя")
			upLb.Move(fyne.NewPos(400, 0))
			dnLb := widget.NewLabel("Нижняя неделя")
			dnLb.Move(fyne.NewPos(400, 215))
			upperSubjectRow.Add(upLb)
			downSubjectRow.Add(dnLb)
			was = true
		}

		upperInitials := subject
		if strings.Contains(subject, " ") {
			upperInitials = utils.GetInitials(subject)
		}
		downInitials := upperInitials

		upperLabel := widget.NewLabel(upperInitials)
		upperLabel.Resize(fyne.NewSize(100, 30))
		upperLabel.Move(fyne.NewPos(0, subjectRowY))
		upperLabel.Wrapping = fyne.TextWrapWord
		upperSubjectRow.Add(upperLabel)

		downLabel := widget.NewLabel(downInitials)
		downLabel.Resize(fyne.NewSize(100, 30))
		downLabel.Move(fyne.NewPos(0, subjectRowY+200))
		downLabel.Wrapping = fyne.TextWrapWord
		downSubjectRow.Add(downLabel)

		dayX := float32(115)
		for _, day := range days {
			upperEntry := widget.NewEntry()
			upperEntry.SetPlaceHolder(day)
			entryWidgets = append(entryWidgets, upperEntry)
			subjectInfo[upperEntry] = subject
			upperEntry.OnChanged = func(text string) {
				updateData(upperEntry, 1)
			}
			upperEntry.Resize(fyne.NewSize(100, 30))
			upperEntry.Move(fyne.NewPos(dayX, subjectRowY))
			upperSubjectRow.Add(upperEntry)

			downEntry := widget.NewEntry()
			downEntry.SetPlaceHolder(day)
			entryWidgets = append(entryWidgets, downEntry)
			subjectInfo[downEntry] = subject
			downEntry.OnChanged = func(text string) {
				updateData(downEntry, 0)
			}
			downEntry.Resize(fyne.NewSize(100, 30))
			downEntry.Move(fyne.NewPos(dayX, subjectRowY+200))
			downSubjectRow.Add(downEntry)

			dayX += float32(115)
		}

		upperW.Add(upperSubjectRow)
		downW.Add(downSubjectRow)

		subjectRowY += float32(35)
	}

	addButton := widget.NewButton("Добавить", func() {
		//for _, entryWidget := range entryWidgets {
		//	updateData(entryWidget)
		//}
	})

	rightContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), addButton, layout.NewSpacer())

	mainContainer := container.New(layout.NewHBoxLayout(),
		container.NewVBox(upperW, downW),
		container.New(layout.NewVBoxLayout(), layout.NewSpacer(), rightContainer))

	return mainContainer
}

func updateData(entry *widget.Entry, week int) {
	if week == 1 {
		subject := subjectInfo[entry]
		fmt.Println("Subject:", subject, "Text:", entry.Text, "WEEk:", "UP")
	} else {
		subject := subjectInfo[entry]
		fmt.Println("Subject:", subject, "Text:", entry.Text, "WEEk:", "DOWN")
	}
}
