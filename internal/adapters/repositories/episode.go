package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type EpisodeRepository interface {
	BulkSave(episodes []entities.Episode) error
}

type GormEpisodeRepository struct {
	db *gorm.DB
}

func NewGormEpisodeRepository(db *gorm.DB) EpisodeRepository {
	return &GormEpisodeRepository{db: db}
}

func (r *GormEpisodeRepository) BulkSave(episodes []entities.Episode) error {
	if result := r.db.Create(&episodes); result.Error != nil {
		return result.Error
	}
	return nil
}
