package initialize

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dgraph-io/badger/v4"
	"log"
	"time"
)

type User struct {
	Name    string    // Имя пользователя
	TimeAdd time.Time // Время добавления
}

func ShowNameInputDialog(w fyne.Window, db *badger.DB) {
	entry := widget.NewEntry()

	lastAddedUser, err := GetLastAddedUserFromBadger(db)
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
				err := addUserToBadger(db, User{Name: userName, TimeAdd: time.Now()})
				if err != nil {
					log.Println("Error saving user name:", err)
					return
				}
			}
		}
	}, w)
}

func GetLastAddedUserFromBadger(db *badger.DB) (User, error) {
	var lastAddedUser User

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()

			if string(key) == models.FileKey {
				continue
			}

			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var user User
			err = utils.Deserialize(value, &user)
			if err != nil {
				return err
			}

			lastAddedUser = user
			break
		}
		return nil
	})
	if err != nil {
		return User{}, err
	}

	return lastAddedUser, nil
}

func addUserToBadger(db *badger.DB, user User) error {
	hashedKey := utils.GenerateHash(user.Name)

	userBytes, err := utils.Serialize(user)
	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(hashedKey), userBytes)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
