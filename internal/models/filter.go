package models

type Filter struct {
	FIO   string
	Group map[string][]map[string][]string //ключ группа ключ 2 тип занятия значение массив дисциплин этой группы
}
