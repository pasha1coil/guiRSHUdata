package repository

import (
	"demofine/internal/models"
	"demofine/internal/utils"
	"github.com/dgraph-io/badger/v4"
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
