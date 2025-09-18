package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeTrailerRepository interface {
	GetByAnimeId(animeId uuid.UUID) ([]entities.AnimeTrailer, error)
	CreateAnimeTrailers(animeTrailers []entities.AnimeTrailer) error
}

type GormAnimeTrailerRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewAnimeTrailerRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeTrailerRepository {
	return &GormAnimeTrailerRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormAnimeTrailerRepository) GetByAnimeId(animeId uuid.UUID) ([]entities.AnimeTrailer, error) {
	var animeTrailers []entities.AnimeTrailer
	if err := r.dbReplica.Where("anime_id = ?", animeId).Find(&animeTrailers).Error; err != nil {
		return nil, err
	}
	return animeTrailers, nil
}

func (r *GormAnimeTrailerRepository) CreateAnimeTrailers(animeTrailers []entities.AnimeTrailer) error {
	if err := r.dbPrimary.Create(&animeTrailers).Error; err != nil {
		return err
	}
	return nil
}
