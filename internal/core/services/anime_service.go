package services

import (
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
)

type AnimeService interface {
	CreateAnime(anime entities.Anime) error
	GetAnimeById(id uint) (*entities.Anime, error)
	GetAnimes(query dtos.AnimeQueryDTO) ([]entities.Anime, error)
	UpdateAnime(anime entities.Anime) error
	DeleteAnime(id uint) error
}

type animeServiceImpl struct {
	repo ports.AnimeRepository
}

func NewAnimeService(repo ports.AnimeRepository) AnimeService {
	return &animeServiceImpl{repo: repo}
}

func (s *animeServiceImpl) CreateAnime(anime entities.Anime) error {
	if err := s.repo.Save(anime); err != nil {
		return err
	}
	return nil
}

func (s *animeServiceImpl) GetAnimeById(id uint) (*entities.Anime, error) {
	anime, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return anime, nil
}

func (s *animeServiceImpl) GetAnimes(query dtos.AnimeQueryDTO) ([]entities.Anime, error) {
	animes, err := s.repo.GetAll(query)
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (s *animeServiceImpl) UpdateAnime(anime entities.Anime) error {
	if err := s.repo.Update(&anime); err != nil {
		return err
	}
	return nil
}

func (s *animeServiceImpl) DeleteAnime(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
