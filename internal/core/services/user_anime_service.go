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
	GetAnimeByUserId(uuid string) ([]dtos.UserAnimeListDTO, error)
	GetMyTopAnime(uuid string) ([]dtos.GetMyTopAnimeResponse, error)
	UpdateMyTopAnime(request *dtos.UpdateMyTopAnimeRequest) error
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

	user, err := s.userRepo.GetByUUID(request.UserUUID)
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

func (s *userAnimeServiceImpl) GetAnimeByUserId(uuid string) ([]dtos.UserAnimeListDTO, error) {
	user, err := s.userRepo.GetByUUID(uuid)
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

func (s *userAnimeServiceImpl) GetMyTopAnime(uuid string) ([]dtos.GetMyTopAnimeResponse, error) {
	user, err := s.userRepo.GetByUUID(uuid)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}

	userAnimes, err := s.repo.GetMyTopAnimeByUserId(user.ID)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	var response []dtos.GetMyTopAnimeResponse
	for _, useranime := range userAnimes {
		response = append(response, dtos.GetMyTopAnimeResponse{
			AnimeID:            useranime.AnimeID,
			AnimeName:          useranime.Anime.Name,
			Score:              useranime.Score,
			Description:        useranime.Anime.Description,
			Episodes:           useranime.Anime.Description,
			Image:              useranime.Anime.Image,
			Status:             useranime.Status,
			WatchAt:            useranime.WatchAt,
			CreatedAt:          useranime.CreatedAt,
			SequenceMyTopAnime: useranime.SequenceMyTopAnime,
		})
	}
	return response, nil
}

func (s *userAnimeServiceImpl) UpdateMyTopAnime(request *dtos.UpdateMyTopAnimeRequest) error {
	user, err := s.userRepo.GetByUUID(request.UserUUID)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewNotFoundError("User not found")
	}

	for _, item := range request.AnimeSequence {
		userAnimes, err := s.repo.GetByUserIdAndAnimeId(user.ID, []uint{item.AnimeID})
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		userAnimeUpdate := userAnimes[0]
		userAnimeUpdate.SequenceMyTopAnime = item.Sequence

		if err := s.repo.UpdateMyTopAnime(&userAnimeUpdate); err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
	}

	return nil
}
