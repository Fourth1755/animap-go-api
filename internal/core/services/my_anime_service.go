package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type MyAnimeService interface {
	AddAnimeToList(userAnime *dtos.AddAnimeToListRequest) error
	GetAnimeByUserId(uuid string) ([]dtos.GetAnimeByUserIdResponse, error)
	GetMyAnimeYearByUserId(uuid string) (*dtos.GetMyAnimeYearByUserIdResponse, error)
	GetMyTopAnime(uuid string) ([]dtos.GetMyTopAnimeResponse, error)
	UpdateMyTopAnime(request *dtos.UpdateMyTopAnimeRequest) error
}

type myAnimeServiceImpl struct {
	repo      repositories.UserAnimeRepository
	animeRepo repositories.AnimeRepository
	userRepo  repositories.UserRepository
}

func NewMyAnimeService(repo repositories.UserAnimeRepository, animeRepo repositories.AnimeRepository, userRepo repositories.UserRepository) MyAnimeService {
	return &myAnimeServiceImpl{repo: repo, animeRepo: animeRepo, userRepo: userRepo}
}

func (s *myAnimeServiceImpl) AddAnimeToList(request *dtos.AddAnimeToListRequest) error {
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
		UserID:        user.ID,
		AnimeID:       request.AnimeID,
		Score:         request.Score,
		Status:        request.Status,
		WatchedYearAt: request.WatchedYear,
	}

	if err := s.repo.Save(&userAnime); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *myAnimeServiceImpl) GetAnimeByUserId(uuid string) ([]dtos.GetAnimeByUserIdResponse, error) {
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

	var animeList []dtos.GetAnimeByUserIdResponse
	for _, useranime := range userAnimes {
		animeList = append(animeList, dtos.GetAnimeByUserIdResponse{
			AnimeID:       useranime.AnimeID,
			AnimeName:     useranime.Anime.Name,
			Score:         useranime.Score,
			Description:   useranime.Anime.Description,
			Episodes:      useranime.Anime.Description,
			Image:         useranime.Anime.Image,
			Status:        useranime.Status,
			WatchedAt:     useranime.WatchedAt,
			WatchedYearAt: useranime.WatchedYearAt,
			CreatedAt:     useranime.CreatedAt,
		})
	}
	return animeList, nil
}

func (s *myAnimeServiceImpl) GetMyAnimeYearByUserId(uuid string) (*dtos.GetMyAnimeYearByUserIdResponse, error) {
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

	var animeYearList []dtos.GetMyAnimeYearByUserIdResponse_AnimeYear
	animeYearMap := make(map[string][]dtos.GetMyAnimeYearByUserIdResponse_AnimeYear_Anime)
	yearList := []string{}
	for _, userAnime := range userAnimes {
		if animeYearMap[userAnime.WatchedYearAt] == nil {
			yearList = append(yearList, userAnime.WatchedYearAt)
		}
		animeYearMap[userAnime.WatchedYearAt] = append(animeYearMap[userAnime.WatchedYearAt], dtos.GetMyAnimeYearByUserIdResponse_AnimeYear_Anime{
			AnimeID:       userAnime.AnimeID,
			AnimeName:     userAnime.Anime.Name,
			Score:         userAnime.Score,
			Description:   userAnime.Anime.Description,
			Episodes:      userAnime.Anime.Description,
			Image:         userAnime.Anime.Image,
			Status:        userAnime.Status,
			WatchedAt:     userAnime.WatchedAt,
			WatchedYearAt: userAnime.WatchedYearAt,
			CreatedAt:     userAnime.CreatedAt,
		})
	}

	for _, year := range yearList {
		if animeYearMap[year] != nil {
			animeYearList = append(animeYearList, dtos.GetMyAnimeYearByUserIdResponse_AnimeYear{
				Year:  year,
				Anime: animeYearMap[year],
			})
		}
	}

	return &dtos.GetMyAnimeYearByUserIdResponse{
		AnimeYear:  animeYearList,
		TotalYear:  uint(len(yearList)),
		TotalAnime: uint(len(userAnimes)),
	}, nil
}

func (s *myAnimeServiceImpl) GetMyTopAnime(uuid string) ([]dtos.GetMyTopAnimeResponse, error) {
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
			WatchAt:            useranime.WatchedAt,
			CreatedAt:          useranime.CreatedAt,
			SequenceMyTopAnime: useranime.SequenceMyTopAnime,
		})
	}
	return response, nil
}

func (s *myAnimeServiceImpl) UpdateMyTopAnime(request *dtos.UpdateMyTopAnimeRequest) error {
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
