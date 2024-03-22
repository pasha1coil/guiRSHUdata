# Генератор ежемесячных отчетов

Приложение для генерации отчетных таблиц на основе нагруочных данных и расписания преподователей.\
`schedule.db` - Папка локальной бд. Тут хранятся данные о пользователях, данные о нагрузке.\
`само приложение` - Стоит сложить его в отдельную папку чтобы бд не потерялась, протянуть к нему ярлых или что то подобное.

## Установка и начало работы
#### 1. Добавьте в `Настройках окружения` ваше имя в том виде как оно прописано в нагрузке пример - `Петров Петр Николаевич`.
#### 2. Добавьте в `Настройках окружения` файл нагрузки за текущий год, принимаются файлы только xslx.

## Логика работы
#### 1.После того как первык пункт прошел успешно отобразятся два меню с названием - `<Месяц> Расписание`. Генерируется всегда только два таких меню, это прошлый и текущий месяцы.
#### 2.После перехода по ним, там уже будут подтянутые группы отфильтрованные по `Группа, Тип занятий` выбираем и `Показать таблицу`.
#### 3.После заполнения данными поля ввода нагрузки по часам на неделю нужно нажать кнопку `Добавить` и данные перейдут в область внутренней памяти приложения.
#### 4.После всех манипуляций с группами и типами занятий лучше стоит убедиться в том что данные все заполнены, если это так, то переходим к `Создать отчет` стоит выбрать путь сохранения файла, если что то забыли ввести, то дозаполнить данные уже не получится.

## Сборки
#### Windows - fyne package -os windows -icon data/assets/logo.png
#### Linux - fyne package -os linux -icon data/assets/logo.png
#### MacOS - fyne package -os darwin -icon data/assets/logo.png
```