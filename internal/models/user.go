package models

import "time"

type User struct {
	Name    string    // Имя пользователя
	TimeAdd time.Time // Время добавления
}
