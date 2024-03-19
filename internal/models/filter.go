package models

type Filter struct {
	FIO   string
	Group map[string][]map[string][]Subjects //ключ группа ключ 2 тип занятия значение массив дисциплин этой группы
}

type Subjects struct {
	Subject string
	Number  string
}
