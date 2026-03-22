package repositories

import (
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeStudioRepository interface {
	Save(animeStudio []entities.AnimeStudio) error
	GetByStudioId(studioId uuid.UUID, cursorTime *time.Time, cursorID *uuid.UUID, limit int) ([]entities.AnimeStudio, error)
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

func (r GormAnimeStudioRepository) GetByStudioId(studioId uuid.UUID, cursorTime *time.Time, cursorID *uuid.UUID, limit int) ([]entities.AnimeStudio, error) {
	var animeStudio []entities.AnimeStudio
	db := r.dbReplica.
		Joins("JOIN animes ON animes.id = anime_studios.anime_id").
		Preload("Anime").
		Where("anime_studios.studio_id = ?", studioId).
		Order("animes.aired_at DESC, animes.id DESC")

	if cursorTime != nil && cursorID != nil {
		db = db.Where("(animes.aired_at < ? OR (animes.aired_at = ? AND animes.id < ?))", *cursorTime, *cursorTime, *cursorID)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}

	result := db.Find(&animeStudio)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeStudio, nil
}
