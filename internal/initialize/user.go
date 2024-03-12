package initialize

import (
	"demofine/internal/models"
	"demofine/internal/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

func ShowNameInputDialog(w fyne.Window, svc *service.Service) {
	entry := widget.NewEntry()

	lastAddedUser, err := svc.Repo.GetLastAddedUserFromBadger()
	if err != nil {
		log.Println("Error retrieving last added user:", err)
		return
	}

	if lastAddedUser.Name != "" {
		entry.SetText(lastAddedUser.Name)
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Name:", entry),
	}

	dialog.ShowForm("Enter Your Name", "Save", "Cancel", formItems, func(accepted bool) {
		if accepted {
			userName := entry.Text
			if userName != "" {
				err := svc.Repo.AddUserToBadger(models.User{Name: userName, TimeAdd: time.Now()})
				if err != nil {
					log.Println("Error saving user name:", err)
					return
				}
			}
		}
	}, w)
}
