package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type UserAnimeService interface {
	AddAnimeToList(userAnime *entities.UserAnime) error
	GetAnimeByUserId(id uint) ([]entities.UserAnime, error)
}

type userAnimeServiceImpl struct {
	repo      repositories.UserAnimeRepository
	animeRepo repositories.AnimeRepository
	userRepo  repositories.UserRepository
}

func NewUserAnimeService(repo repositories.UserAnimeRepository, animeRepo repositories.AnimeRepository, userRepo repositories.UserRepository) UserAnimeService {
	return &userAnimeServiceImpl{repo: repo, animeRepo: animeRepo, userRepo: userRepo}
}

func (s *userAnimeServiceImpl) AddAnimeToList(userAnime *entities.UserAnime) error {
	if _, err := s.animeRepo.GetById(userAnime.AnimeID); err != nil {
		logs.Error(err.Error())
		return err
	}

	if _, err := s.userRepo.GetById(userAnime.UserID); err != nil {
		logs.Error(err.Error())
		return err
	}

	if err := s.repo.Save(userAnime); err != nil {
		logs.Error(err.Error())
		return err
	}

	return nil
}

func (s *userAnimeServiceImpl) GetAnimeByUserId(id uint) ([]entities.UserAnime, error) {
	if _, err := s.userRepo.GetById(id); err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	userAnimes, err := s.repo.GetByUserId(id)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return userAnimes, nil
}
