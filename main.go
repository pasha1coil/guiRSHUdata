package main

import (
	"demofine/data"
	"demofine/internal/initialize"
	"demofine/internal/models"
	"demofine/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/dgraph-io/badger/v4"
	"log"
	"sync"
)

var (
	daysMutex sync.Mutex
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		days := utils.GenerateDays()

		daysMutex.Lock()
		defer daysMutex.Unlock()

		models.DaysInfo = days
	}()

	go func() {
		wg.Wait()
	}()

	a := app.NewWithID("RSHU.reports")
	a.SetIcon(data.RSHULogo)
	initialize.LogLifecycle(a)
	db, err := badger.Open(badger.DefaultOptions(models.DbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	w := a.NewWindow("RSHU reports")
	w.SetMaster()

	models.TopWindow = w
	menu := initialize.MakeMenu(a, w, db)
	if menu == nil {
		panic("menu is nil")
	}
	w.SetMainMenu(menu)
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}
