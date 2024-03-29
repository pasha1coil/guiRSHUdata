package models

import (
	"fyne.io/fyne/v2"
	"time"
)

type Table struct {
	Title string
	View  func(w fyne.Window, month string) fyne.CanvasObject
	Month string
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
	Month    string
	Subject  string
	Group    string
	Type     string
	Number   string
	UpperDay map[string]map[string][1]Subjects
	LowerDay map[string]map[string][1]Subjects
}

type GenerateReport struct {
	Month    string
	Group    string
	Type     string
	Subject  string
	DayWeek  string
	TypeWeek WeekType
	Number   string
	Entry    string
	Created  time.Time
}
