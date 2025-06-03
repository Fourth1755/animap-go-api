package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeService interface {
	CreateEpisode(id uuid.UUID) error
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
		episodes = append(episodes, entities.Episode{
			AnimeID: anime_id,
			Number:  uint(i),
		})
	}

	err = s.episodeRepo.BulkSave(episodes)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}
