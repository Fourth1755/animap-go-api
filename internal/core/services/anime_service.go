package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
)

type AnimeService interface {
	CreateAnime(anime entities.Anime) error
	GetAnimeById(id uint) (*entities.Anime, error)
	GetAnimes(query dtos.AnimeQueryDTO) ([]dtos.AnimeDTO, error)
	UpdateAnime(anime entities.Anime) error
	DeleteAnime(id uint) error
	GetAnimeByUserId(user_id uint) ([]entities.UserAnime, error)
}

type animeServiceImpl struct {
	repo     repositories.AnimeRepository
	userRepo repositories.UserRepository
}

func NewAnimeService(repo repositories.AnimeRepository, userRepo repositories.UserRepository) AnimeService {
	return &animeServiceImpl{repo: repo, userRepo: userRepo}
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

func (s *animeServiceImpl) GetAnimes(query dtos.AnimeQueryDTO) ([]dtos.AnimeDTO, error) {
	animes, err := s.repo.GetAll(query)
	if err != nil {
		return nil, err
	}

	var animesDto []dtos.AnimeDTO
	for _, anime := range animes {
		animesDto = append(animesDto, dtos.AnimeDTO{
			ID:       anime.ID,
			Name:     anime.Name,
			Episodes: anime.Episodes,
			Seasonal: anime.Seasonal,
			Year:     anime.Year,
		})
	}
	return animesDto, nil
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

func (s *animeServiceImpl) GetAnimeByUserId(user_id uint) ([]entities.UserAnime, error) {
	if _, err := s.userRepo.GetById(user_id); err != nil {
		return nil, err
	}
	result, err := s.repo.GetByUserId(user_id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
