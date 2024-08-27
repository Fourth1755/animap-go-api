package services

import (
	"gorm.io/gorm"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type AnimeService interface {
	CreateAnime(anime entities.Anime) error
	GetAnimeById(id uint) (*dtos.AnimeDetailResponse, error)
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
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *animeServiceImpl) GetAnimeById(id uint) (*dtos.AnimeDetailResponse, error) {
	anime, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("Anime not found")
		}
		return nil, errs.NewUnexpectedError()
	}

	var categories []dtos.AnimeDetailCategories
	for _, category := range anime.Categories {
		categories = append(categories, dtos.AnimeDetailCategories{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	animeResponse := dtos.AnimeDetailResponse{
		ID:          anime.ID,
		Name:        anime.Name,
		NameEnglish: anime.NameEnglish,
		Episodes:    anime.Episodes,
		Seasonal:    anime.Seasonal,
		Year:        anime.Year,
		Image:       anime.Image,
		Description: anime.Description,
		Type:        anime.Type,
		Duration:    anime.Duration,
		Categories:  categories,
	}
	return &animeResponse, nil
}

func (s *animeServiceImpl) GetAnimes(query dtos.AnimeQueryDTO) ([]dtos.AnimeDTO, error) {
	animes, err := s.repo.GetAll(query)
	if err != nil {
		logs.Error(err.Error())
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
	_, err := s.repo.GetById(anime.ID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}

	if err := s.repo.Update(&anime); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *animeServiceImpl) DeleteAnime(id uint) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}

	if err := s.repo.Delete(id); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *animeServiceImpl) GetAnimeByUserId(user_id uint) ([]entities.UserAnime, error) {
	if _, err := s.userRepo.GetById(user_id); err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}
	result, err := s.repo.GetByUserId(user_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	return result, nil
}
