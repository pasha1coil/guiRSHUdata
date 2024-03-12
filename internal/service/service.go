package service

import (
	"demofine/internal/models"
	"demofine/internal/repository"
	"demofine/internal/utils"
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
	title := month + " Расписание"
	models.Tables[month] = models.Table{
		Title: title,
		View:  s.MakeTableTab,
	}
	models.TableIndex[""] = append(models.TableIndex[""], month)
}