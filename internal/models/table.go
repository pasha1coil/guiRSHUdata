package models

import "fyne.io/fyne/v2"

type Table struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

var (
	Tables     = map[string]Table{}
	TableIndex = map[string][]string{}
)

// так как тип непонятен используем any
type LoadedFile struct {
	Number         []any
	Plan           []any
	Faculty        []any
	Block          []any
	Discipline     []any
	Semester       []any
	Group          []any
	CountStudents  []any
	Type           []any
	Hours          []any
	FIO            []any
	TypeWork       []any
	RankClever     []any
	ProgressClever []any
}

type EntryData struct {
	Subject  string
	Group    string
	Type     string
	Number   string
	UpperDay map[string]map[Subjects]string
	LowerDay map[string]map[Subjects]string
}
