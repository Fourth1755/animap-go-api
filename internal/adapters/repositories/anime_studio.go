package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type AnimeStudioRepository interface {
	Save(animeStudio []entities.AnimeStudio) error
}

type GormAnimeStudioRepository struct {
	db *gorm.DB
}

func NewGormAnimeStudioRepository(db *gorm.DB) AnimeStudioRepository {
	return &GormAnimeStudioRepository{db: db}
}

func (r GormAnimeStudioRepository) Save(animeStudio []entities.AnimeStudio) error {
	if result := r.db.Create(&animeStudio); result.Error != nil {
		return result.Error
	}
	return nil
}
