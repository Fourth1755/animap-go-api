package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeCharacterRepository interface {
	BulkSave(episodeCharacter []entities.EpisodeCharacter) error
	GetByEpisodeIds(episode_ids []uuid.UUID) ([]entities.EpisodeCharacter, error)
}

type GormEpisodeCharacterRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormEpisodeCharacterRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) EpisodeCharacterRepository {
	return &GormEpisodeCharacterRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormEpisodeCharacterRepository) BulkSave(episodeCharacter []entities.EpisodeCharacter) error {
	if result := r.dbPrimary.Create(&episodeCharacter); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormEpisodeCharacterRepository) GetByEpisodeIds(episode_ids []uuid.UUID) ([]entities.EpisodeCharacter, error) {
	var episodeCharacters []entities.EpisodeCharacter
	if result := r.dbReplica.
		Preload("Character").
		Where("episode_id in (?)", episode_ids).Find(&episodeCharacters); result.Error != nil {
		return nil, result.Error
	}
	return episodeCharacters, nil
}
