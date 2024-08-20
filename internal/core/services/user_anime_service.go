package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type UserAnimeService interface {
	AddAnimeToList(userAnime *entities.UserAnime) error
	GetAnimeByUserId(id uint) ([]dtos.UserAnimeListDTO, error)
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
		return errs.NewNotFoundError("Anime not found")
	}

	if _, err := s.userRepo.GetById(userAnime.UserID); err != nil {
		logs.Error(err.Error())
		return errs.NewNotFoundError("User not found")
	}

	if err := s.repo.Save(userAnime); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *userAnimeServiceImpl) GetAnimeByUserId(id uint) ([]dtos.UserAnimeListDTO, error) {
	if _, err := s.userRepo.GetById(id); err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}

	userAnimes, err := s.repo.GetByUserId(id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	var animeList []dtos.UserAnimeListDTO
	for _, useranime := range userAnimes {
		animeList = append(animeList, dtos.UserAnimeListDTO{
			AnimeID:     useranime.AnimeID,
			AnimeName:   useranime.Anime.Name,
			Score:       useranime.Score,
			Description: useranime.Anime.Description,
			Episodes:    useranime.Anime.Description,
			Image:       useranime.Anime.Image,
			Status:      useranime.Status,
			WatchAt:     useranime.WatchAt,
			CreatedAt:   useranime.CreatedAt,
		})
	}
	return animeList, nil
}
