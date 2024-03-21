package repository

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"github.com/dgraph-io/badger/v4"
	"strings"
)

type Repository struct {
	Db *badger.DB
}

func NewRepository(db *badger.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

func (r *Repository) GetLastAddedUserFromBadger() (models.User, error) {
	var lastAddedUser models.User

	err := r.Db.View(func(txn *badger.Txn) error {
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

			var user models.User
			err = utils.Deserialize(value, &user)
			if err != nil {
				return err
			}

			if lastAddedUser.TimeAdd.Unix() < user.TimeAdd.Unix() {
				lastAddedUser = user
			}
		}
		return nil
	})
	if err != nil {
		return models.User{}, err
	}

	return lastAddedUser, nil
}

func (r *Repository) AddUserToBadger(user models.User) error {
	hashedKey := utils.GenerateHash(user.Name)

	userBytes, err := utils.Serialize(user)
	if err != nil {
		return err
	}

	err = r.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(hashedKey), userBytes)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckFileLoaded() (error, bool) {
	var checker bool
	err := r.Db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(models.FileKey))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				checker = false
				return nil
			}
			checker = false
			return err
		}
		checker = true
		return nil
	})
	if err != nil {
		return err, checker
	}
	return nil, checker
}

func (r *Repository) ReadFileFromBadger() ([]byte, error) {
	var data []byte
	err := r.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(models.FileKey))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Repository) AddFileToBadger(fileHash string, fileData []byte) error {
	err := r.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(fileHash), fileData)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFileData(fileName string) ([]byte, error) {
	fileHash := utils.GenerateHash(fileName)
	var fileData []byte
	err := r.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fileHash))
		if err != nil {
			return err
		}
		fileData, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func (r *Repository) GetAllFileNames() ([]string, error) {
	var fileNames []string
	err := r.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			fileName := string(item.Key())
			if strings.HasPrefix(fileName, "report") {
				fileNames = append(fileNames, fileName)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}
