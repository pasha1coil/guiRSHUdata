package main

import (
	"demofine/data"
	"demofine/internal/initialize"
	"demofine/internal/models"
	"fyne.io/fyne/v2/app"
	"github.com/dgraph-io/badger/v4"
	"log"
)

func main() {
	a := app.NewWithID("RSHU.demo")
	a.SetIcon(data.FyneLogo)
	initialize.LogLifecycle(a)
	w := a.NewWindow("Fyne Demo")
	models.TopWindow = w
	db, err := badger.Open(badger.DefaultOptions(models.DbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	w.SetMainMenu(initialize.MakeMenu(a, w, db))
	w.SetMaster()
	w.ShowAndRun()
}
