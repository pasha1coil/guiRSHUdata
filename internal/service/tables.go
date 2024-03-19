package service

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/tealeg/xlsx"
	"log"
	"strings"
)

var days = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}

var FinishData = map[string][]models.EntryData{}

func (s *Service) MakeTableTab(w fyne.Window, month string) fyne.CanvasObject {

	FinishData = make(map[string][]models.EntryData)

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
	filter.Group = make(map[string][]map[string][]models.Subjects)
	for i, teacher := range loadedFile.FIO {
		if teacher == userName {
			group := loadedFile.Group[i].(string)
			discipline := loadedFile.Discipline[i].(string)
			lessonType := loadedFile.Type[i].(string)
			docNumber := loadedFile.Number[i].(string)
			if _, ok := filter.Group[group]; !ok {
				filter.Group[group] = make([]map[string][]models.Subjects, 0)
			}

			var found bool
			for j := range filter.Group[group] {
				if filter.Group[group][j][lessonType] != nil {
					filter.Group[group][j][lessonType] = append(filter.Group[group][j][lessonType], models.Subjects{Subject: discipline, Number: docNumber})
					found = true
				}
			}
			if !found {
				newLesson := make(map[string][]models.Subjects)
				newLesson[lessonType] = append(newLesson[lessonType], models.Subjects{Subject: discipline, Number: docNumber})
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

	content := container.NewVBox(
		widget.NewLabel("Вы уверены, что хотите продолжить?"),
	)

	createReportButton := widget.NewButton("Создать отчет", func() {
		confirm := dialog.NewCustomConfirm("Вы уверены, что хотите создать отчет?", "ДА", "Отмена", content, func(confirmed bool) {
			if confirmed {
				s.generateReport(FinishData)
			}
		}, w)

		confirm.Show()
	})

	dialog := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		selectTab,
		selectType,
		widget.NewButton("Показать таблицу", func() {
			var data []models.Subjects
			var ok bool
			for _, dis := range filter.Group[selectTab.Selected] {
				if data, ok = dis[selectType.Selected]; ok {
					break
				}
			}

			tableWindow := fyne.CurrentApp().NewWindow("Таблица")
			table := createTable(data, selectTab.Selected, selectType.Selected, month, tableWindow)
			tableWindow.SetContent(table)
			tableWindow.Resize(fyne.NewSize(800, 600))
			tableWindow.Show()
		}),
		createReportButton,
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

func createTable(subjects []models.Subjects, group string, lessonType string, month string, tableWindow fyne.Window) fyne.CanvasObject {
	var subjectInfoUp = make(map[*widget.Entry]models.Subjects)
	var subjectInfoDown = make(map[*widget.Entry]models.Subjects)
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

		upperInitials := subject.Subject
		if strings.Contains(subject.Subject, " ") {
			upperInitials = utils.GetInitials(subject.Subject)
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
			subjectInfoUp[upperEntry] = subject
			upperEntry.OnChanged = func(text string) {
				updateData(upperEntry, 1, group, lessonType, subjectInfoUp, month)
			}
			upperEntry.Resize(fyne.NewSize(100, 30))
			upperEntry.Move(fyne.NewPos(dayX, subjectRowY))
			upperSubjectRow.Add(upperEntry)

			downEntry := widget.NewEntry()
			downEntry.SetPlaceHolder(day)
			entryWidgets = append(entryWidgets, downEntry)
			subjectInfoDown[downEntry] = subject
			downEntry.OnChanged = func(text string) {
				updateData(downEntry, 0, group, lessonType, subjectInfoDown, month)
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
		tableWindow.Close()
	})

	rightContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), addButton, layout.NewSpacer())

	mainContainer := container.New(layout.NewHBoxLayout(),
		container.NewVBox(upperW, downW),
		container.New(layout.NewVBoxLayout(), layout.NewSpacer(), rightContainer))

	return mainContainer
}

func updateData(entry *widget.Entry, week int, group string, lessonType string, subjectInfo map[*widget.Entry]models.Subjects, month string) {
	subject := subjectInfo[entry]
	key := utils.GenerateHash(group + lessonType)

	if FinishData[key] == nil {
		FinishData[key] = make([]models.EntryData, 0)
	}

	var entryData *models.EntryData
	for i, ed := range FinishData[key] {
		if ed.Subject == subject.Subject {
			entryData = &FinishData[key][i]
			break
		}
	}
	if entryData == nil {
		entryData = &models.EntryData{
			Month:    month,
			Subject:  subject.Subject,
			Group:    group,
			Number:   subject.Number,
			Type:     lessonType,
			UpperDay: make(map[string]map[models.Subjects]string),
			LowerDay: make(map[string]map[models.Subjects]string),
		}
		FinishData[key] = append(FinishData[key], *entryData)
	}

	dayMap := entryData.UpperDay
	if week == 0 {
		dayMap = entryData.LowerDay
	}

	day := entry.PlaceHolder
	if dayMap[day] == nil {
		dayMap[day] = make(map[models.Subjects]string)
	}
	dayMap[day][subject] = entry.Text
}
