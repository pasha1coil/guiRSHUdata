package initialize

import (
	"demofine/internal/models"
	"demofine/internal/repository"
	"demofine/internal/service"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dgraph-io/badger/v4"
	"io/ioutil"
	"log"
	"strings"
)

func MakeMenu(a fyne.App, w fyne.Window, db *badger.DB) *fyne.MainMenu {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	svc.InstallTables()

	nameInputDialogItem := fyne.NewMenuItem("Enter Name", func() {
		ShowNameInputDialog(w, svc)
	})

	loadScheduleItem := fyne.NewMenuItem("Load Schedule", func() {
		done := make(chan struct{})
		go loadSchedule(db, done)
		go func() {
			<-done
			showRemainingMenuItems(a, w, svc)
		}()
	})

	nameInputDialogItem.Icon = theme.DocumentCreateIcon()
	loadScheduleItem.Icon = theme.ContentPasteIcon()

	file := fyne.NewMenu("File", nil)
	main := fyne.NewMainMenu(file)

	file.Items = []*fyne.MenuItem{nameInputDialogItem, loadScheduleItem}

	showRemainingMenuItems(a, w, svc)

	return main
}

func showRemainingMenuItems(a fyne.App, w fyne.Window, svc *service.Service) {
	err, checker := svc.Repo.CheckFileLoaded()
	if err != nil {
		log.Println("Error checking file loaded status:", err)
		return
	}

	if checker {

		content := container.NewStack()
		title := widget.NewLabel("Component name")
		setTutorial := func(t models.Table) {
			if fyne.CurrentDevice().IsMobile() {
				child := a.NewWindow(t.Title)
				models.TopWindow = child
				child.SetContent(t.View(models.TopWindow))
				child.Show()
				child.SetOnClosed(func() {
					models.TopWindow = w
				})
				return
			}

			title.SetText(t.Title)

			content.Objects = []fyne.CanvasObject{t.View(w)}
			content.Refresh()
		}

		tutorial := container.NewBorder(
			container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
		if err == nil {
			if fyne.CurrentDevice().IsMobile() {
				w.SetContent(makeNav(setTutorial, false))
			} else {
				split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
				split.Offset = 0.2
				w.SetContent(split)
			}
		}
	}
}

func loadSchedule(db *badger.DB, done chan struct{}) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		defer func() {
			done <- struct{}{}
		}()

		if err != nil {
			dialog.ShowError(err, models.TopWindow)
			return
		}

		if reader == nil {
			return
		}

		defer reader.Close()

		fileName := reader.URI().Path()
		if !strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
			dialog.ShowError(errors.New("only Excel files (.xlsx) are allowed"), models.TopWindow)
			return
		}

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, models.TopWindow)
			return
		}

		err = db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(models.FileKey), data)
			return err
		})
		if err != nil {
			dialog.ShowError(err, models.TopWindow)
			return
		}

		dialog.ShowInformation("File Loaded", "Schedule file successfully loaded", models.TopWindow)
	}, models.TopWindow)
}

func makeNav(setTutorial func(tutorial models.Table), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return models.TableIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := models.TableIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := models.Tables[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := models.Tables[uid]; ok {
				a.Preferences().SetString(models.PreferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(models.PreferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
