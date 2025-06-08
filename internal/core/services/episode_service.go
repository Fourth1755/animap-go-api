package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeService interface {
	CreateEpisode(id uuid.UUID) error
	GetByAnimeId(anime_id uuid.UUID) (*dtos.GetEpisodeResponse, error)
	UpdateEpisode(request dtos.UpdateEpisodeRequest) error
}

type episodeServiceImpl struct {
	episodeRepo repositories.EpisodeRepository
	animeRepo   repositories.AnimeRepository
}

func NewEpisodeService(episodeRepo repositories.EpisodeRepository, animeRepo repositories.AnimeRepository) EpisodeService {
	return &episodeServiceImpl{
		episodeRepo: episodeRepo,
		animeRepo:   animeRepo}
}

func (s *episodeServiceImpl) CreateEpisode(animeId uuid.UUID) error {
	anime, err := s.animeRepo.GetById(animeId)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}
	if anime.IsCreateEpisode {
		return errs.NewBadRequestError("This Anime had already been create episode.")
	}

	var episodes []entities.Episode
	for i := 1; i <= anime.Episodes; i++ {
		episodeId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		episodes = append(episodes, entities.Episode{
			ID:      episodeId,
			AnimeID: animeId,
			Number:  uint(i),
		})
	}

	err = s.episodeRepo.BulkSave(episodes)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	err = s.animeRepo.UpdateIsCreateEpisode(anime.ID)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *episodeServiceImpl) GetByAnimeId(anime_id uuid.UUID) (*dtos.GetEpisodeResponse, error) {
	episodes, err := s.episodeRepo.GetByAnimeId(anime_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	var episodeResponse []dtos.GetEpisodeResponseEpisode
	for _, episode := range episodes {
		episodeResponse = append(episodeResponse, dtos.GetEpisodeResponseEpisode{
			ID:          episode.ID,
			Number:      episode.Number,
			Name:        episode.Name,
			NameThai:    episode.NameThai,
			NameEnglish: episode.NameEnglish,
		})
	}
	return &dtos.GetEpisodeResponse{
		Episodes: episodeResponse,
	}, nil
}

func (s *episodeServiceImpl) UpdateEpisode(request dtos.UpdateEpisodeRequest) error {
	_, err := s.episodeRepo.GetById(request.ID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Episode not found")
		}
		return errs.NewUnexpectedError()
	}

	episodeUpdate := entities.Episode{
		ID:          request.ID,
		Name:        request.Name,
		NameThai:    request.NameThai,
		NameEnglish: request.NameEnglish,
	}

	if err := s.episodeRepo.Update(&episodeUpdate); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}
