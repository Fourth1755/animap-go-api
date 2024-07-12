package services

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
)

type UserAnimeService interface {
	AddAnimeToList(userAnime *entities.UserAnime) error
	GetAnimeByUserId(id uint) ([]entities.UserAnime, error)
}

type userAnimeServiceImpl struct {
	repo      ports.UserAnimeRepository
	animeRepo ports.AnimeRepository
	userRepo  ports.UserRepository
}

func NewUserAnimeService(repo ports.UserAnimeRepository, animeRepo ports.AnimeRepository, userRepo ports.UserRepository) UserAnimeService {
	return &userAnimeServiceImpl{repo: repo, animeRepo: animeRepo, userRepo: userRepo}
}

func (s *userAnimeServiceImpl) AddAnimeToList(userAnime *entities.UserAnime) error {
	if _, err := s.animeRepo.GetById(userAnime.AnimeID); err != nil {
		return err
	}

	if _, err := s.userRepo.GetById(userAnime.UserID); err != nil {
		return err
	}

	if err := s.repo.Save(userAnime); err != nil {
		return err
	}

	return nil
}

func (s *userAnimeServiceImpl) GetAnimeByUserId(id uint) ([]entities.UserAnime, error) {
	if _, err := s.userRepo.GetById(id); err != nil {
		return nil, err
	}

	userAnimes, err := s.repo.GetByUserId(id)
	if err != nil {
		return nil, err
	}

	return userAnimes, nil
}
