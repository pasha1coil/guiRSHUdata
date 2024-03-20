package main

import (
	"demofine/data"
	"demofine/internal/initialize"
	"demofine/internal/models"
	"demofine/internal/service"
	"fyne.io/fyne/v2/app"
	"github.com/dgraph-io/badger/v4"
	"log"
)

func main() {
	models.Dates = service.GenerateDays()

	a := app.NewWithID("RSHU.demo")
	a.SetIcon(data.FyneLogo)
	initialize.LogLifecycle(a)
	db, err := badger.Open(badger.DefaultOptions(models.DbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	w := a.NewWindow("RSHU Demo")
	w.SetMaster()

	models.TopWindow = w
	menu := initialize.MakeMenu(a, w, db)
	if menu == nil {
		panic("menu is nil")
	}
	w.SetMainMenu(menu)
	w.ShowAndRun()
}
