package tutorials

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var headers = []string{
	"№ п/п", "Название предмета", "Дата", "Факультет, курс, группа", "ТЕМА ЗАНЯТИЙ",
	"Очная и очно-заочная формы обучения", "Заочная форма обучения", "Лекция", "Практические",
	"Лабораторные работы", "Консультации", "Зачет", "Экзамен", "Учебная практика",
	"Производственная практика", "Преддипломная практика", "Курсовые работы", "Руководство ВКР",
	"ГИА", "Руководство аспирантами", "Занятия с аспирантами",
}

func makeTableTab(_ fyne.Window) fyne.CanvasObject {
	var entryGrid [][]*widget.Entry

	var headerLabels []fyne.CanvasObject
	for _, header := range headers {
		label := widget.NewLabel(header)
		headerLabels = append(headerLabels, label)
	}

	headersColumn := container.New(layout.NewGridLayoutWithRows(len(headers)), headerLabels...)

	leftColumn := container.New(layout.NewGridLayoutWithRows(len(headers)))

	entryGrid = make([][]*widget.Entry, len(headers))
	for i := range entryGrid {
		entryGrid[i] = make([]*widget.Entry, 0)
	}

	for range headers {
		var entryLabels []fyne.CanvasObject
		for range headers {
			entry := widget.NewEntry()
			entry.MultiLine = true
			entryLabels = append(entryLabels, entry)
		}
		leftColumn.Add(container.New(layout.NewGridLayoutWithColumns(len(entryLabels)), entryLabels...))
	}

	headersScroll := container.NewScroll(headersColumn)
	headersScroll.SetMinSize(fyne.NewSize(200, 200))

	entriesScroll := container.NewHScroll(leftColumn)
	entriesScroll.SetMinSize(fyne.NewSize(600, 200))

	addColumnButton := widget.NewButton("Добавить", func() {
		addColumn(entryGrid, leftColumn)
	})

	addColumnButton.Resize(fyne.NewSize(100, 30))

	content := container.New(layout.NewHBoxLayout(),
		headersScroll,
		entriesScroll,
		container.NewVBox(layout.NewSpacer(), addColumnButton),
	)

	return content
}

func addColumn(entryGrid [][]*widget.Entry, leftColumn *fyne.Container) {
	for i := range entryGrid {
		entry := widget.NewEntry()
		entry.MultiLine = true
		entryGrid[i] = append(entryGrid[i], entry)
		leftColumn.Objects[i].(*fyne.Container).Add(entry)
	}
}
