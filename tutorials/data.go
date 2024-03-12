package tutorials

import (
	"fyne.io/fyne/v2"
)

type Table struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

var (
	Tables = map[string]Table{}

	TableIndex = map[string][]string{}
)

func init() {
	months := []string{"September", "October", "November", "December", "January", "February", "March", "April", "May", "June", "July", "August"}
	for _, month := range months {
		title := month + " Расписание"
		Tables[month] = Table{
			Title: title,
			View:  makeTableTab,
		}
		TableIndex[""] = append(TableIndex[""], month)
	}
}
