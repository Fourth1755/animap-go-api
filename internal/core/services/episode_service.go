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

func (s *episodeServiceImpl) CreateEpisode(anime_id uuid.UUID) error {
	anime, err := s.animeRepo.GetById(anime_id)
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
			AnimeID: anime_id,
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
