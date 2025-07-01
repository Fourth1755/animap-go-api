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
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeStudioRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeStudioRepository {
	return &GormAnimeStudioRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormAnimeStudioRepository) Save(animeStudio []entities.AnimeStudio) error {
	if result := r.dbPrimary.Create(&animeStudio); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormAnimeStudioRepository) GetByStudioId(studioId uuid.UUID) ([]entities.AnimeStudio, error) {
	var animeStudio []entities.AnimeStudio
	result := r.dbReplica.Preload("Anime").Where("studio_id = ?", studioId).Find(&animeStudio)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeStudio, nil
}
