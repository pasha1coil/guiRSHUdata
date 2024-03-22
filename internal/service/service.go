package service

import (
	"demofine/internal/models"
	"demofine/internal/repository"
	"demofine/internal/utils"
	"fyne.io/fyne/v2"
)

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) LoadDataFromBadger() ([]byte, string, error) {
	fileData, err := s.Repo.ReadFileFromBadger()
	if err != nil {
		return nil, "", err
	}

	lastUser, err := s.Repo.GetLastAddedUserFromBadger()
	if err != nil {
		return nil, "", err
	}

	return fileData, lastUser.Name, nil
}

func (s *Service) InstallTables() {
	currentMonth := utils.GetCurrentMonth()
	previousMonth := utils.GetPreviousMonth()

	s.createTable(currentMonth)
	s.createTable(previousMonth)
}

func (s *Service) createTable(month string) {
	title := models.RussianMonth[month] + " Расписание"
	viewFunc := func(w fyne.Window, month string) fyne.CanvasObject {
		return s.MakeTableTab(w, month)
	}

	models.Tables[month] = models.Table{
		Title: title,
		View:  viewFunc,
		Month: month,
	}
	models.TableIndex[""] = append(models.TableIndex[""], month)
}

func (s *Service) Help() {
	viewFunc := func(w fyne.Window, month string) fyne.CanvasObject {
		return s.WelcomeScreen(w, month)
	}
	models.Tables["Документация"] = models.Table{
		Title: "Документация",
		View:  viewFunc,
		Month: "Документация",
	}
	models.TableIndex[""] = append(models.TableIndex[""], "Документация")
}
