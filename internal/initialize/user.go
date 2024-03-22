package initialize

import (
	"demofine/internal/models"
	"demofine/internal/service"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"time"
)

func ShowNameInputDialog(w fyne.Window, svc *service.Service) {
	entry := widget.NewEntry()

	lastAddedUser, err := svc.Repo.GetLastAddedUserFromBadger()
	if err != nil {
		errorMessage := "Ошибка при получении последнего добавленного пользователя: " + err.Error()
		dialog.ShowError(errors.New(errorMessage), models.TopWindow)
	}

	if lastAddedUser.Name != "" {
		entry.SetText(lastAddedUser.Name)
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Name:", entry),
	}

	dialog.ShowForm("Введи ваше имя", "Сохранить", "Отмена", formItems, func(accepted bool) {
		if accepted {
			userName := entry.Text
			if userName != "" {
				err := svc.Repo.AddUserToBadger(models.User{Name: userName, TimeAdd: time.Now()})
				if err != nil {
					errorMessage := "Ошибка сохранения имени пользователя: " + err.Error()
					dialog.ShowError(errors.New(errorMessage), models.TopWindow)
				}
			}
		}
	}, w)
}
