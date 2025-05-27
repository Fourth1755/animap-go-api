package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeStudioRepository interface {
	Save(animeStudio []entities.AnimeStudio) error
	GetByStudioId(studioId uuid.UUID) ([]entities.AnimeStudio, error)
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

func (r GormAnimeStudioRepository) GetByStudioId(studioId uuid.UUID) ([]entities.AnimeStudio, error) {
	var animeStudio []entities.AnimeStudio
	result := r.db.Preload("Anime").Where("studio_id = ?", studioId).Find(&animeStudio)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeStudio, nil
}
