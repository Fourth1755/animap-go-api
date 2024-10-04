package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type UserAnimeService interface {
	AddAnimeToList(userAnime *dtos.AddAnimeToListRequest) error
	GetAnimeByUserId(sid string) ([]dtos.UserAnimeListDTO, error)
}

type userAnimeServiceImpl struct {
	repo      repositories.UserAnimeRepository
	animeRepo repositories.AnimeRepository
	userRepo  repositories.UserRepository
}

func NewUserAnimeService(repo repositories.UserAnimeRepository, animeRepo repositories.AnimeRepository, userRepo repositories.UserRepository) UserAnimeService {
	return &userAnimeServiceImpl{repo: repo, animeRepo: animeRepo, userRepo: userRepo}
}

func (s *userAnimeServiceImpl) AddAnimeToList(request *dtos.AddAnimeToListRequest) error {
	if _, err := s.animeRepo.GetById(request.AnimeID); err != nil {
		logs.Error(err.Error())
		return errs.NewNotFoundError("Anime not found")
	}

	user, err := s.userRepo.GetBySid(request.Sid)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewNotFoundError("User not found")
	}
	userAnime := entities.UserAnime{
		UserID:  user.ID,
		AnimeID: request.AnimeID,
		Score:   request.Score,
		Status:  request.Status,
	}

	if err := s.repo.Save(&userAnime); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *userAnimeServiceImpl) GetAnimeByUserId(sid string) ([]dtos.UserAnimeListDTO, error) {
	user, err := s.userRepo.GetBySid(sid)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}

	userAnimes, err := s.repo.GetByUserId(user.ID)
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
