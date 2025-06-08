package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type EpisodeCharacterRepository interface {
	BulkSave(episodeCharacter []entities.EpisodeCharacter) error
}

type GormEpisodeCharacterRepository struct {
	db *gorm.DB
}

func NewGormEpisodeCharacterRepository(db *gorm.DB) EpisodeCharacterRepository {
	return &GormEpisodeCharacterRepository{db: db}
}

func (r *GormEpisodeCharacterRepository) BulkSave(episodeCharacter []entities.EpisodeCharacter) error {
	if result := r.db.Create(&episodeCharacter); result.Error != nil {
		return result.Error
	}
	return nil
}
