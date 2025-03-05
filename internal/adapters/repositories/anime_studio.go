package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type AnimeStudioRepository interface {
	GetAll() ([]entities.Studio, error)
}

type GormAnimeStudioRepository struct {
	db *gorm.DB
}

func NewGormAnimeStudioRepository(db *gorm.DB) SongRepository {
	return &GormSongRepository{db: db}
}

func (r GormAnimeStudioRepository) Save(animeStudio []entities.AnimeStudio) error {
	if result := r.db.Create(&animeStudio); result.Error != nil {
		return result.Error
	}
	return nil
}
