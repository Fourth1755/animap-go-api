package repositories

import (
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCategoryUniverseRepository interface {
	Save(animeCategory []entities.AnimeCategoryUniverse) error
	GetByCategoryUniverseId(categoryID uuid.UUID, cursorTime *time.Time, cursorID *uuid.UUID, limit int) ([]entities.AnimeCategoryUniverse, error)
	GetByAnimeIdsAndCategoryUniverseIds(anime_ids []uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategoryUniverse, error)
}

type GormAnimeCategoryUniverseRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeCategoryUniverseRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeCategoryUniverseRepository {
	return &GormAnimeCategoryUniverseRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormAnimeCategoryUniverseRepository) Save(animeCategory []entities.AnimeCategoryUniverse) error {
	if result := r.dbPrimary.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormAnimeCategoryUniverseRepository) GetByCategoryUniverseId(categoryID uuid.UUID, cursorTime *time.Time, cursorID *uuid.UUID, limit int) ([]entities.AnimeCategoryUniverse, error) {
	var categoryAnime []entities.AnimeCategoryUniverse
	db := r.dbReplica.
		Joins("JOIN animes ON animes.id = anime_category_universes.anime_id").
		Preload("Anime").
		Where("anime_category_universes.category_universe_id = ?", categoryID).
		Order("animes.aired_at DESC, animes.id DESC")

	if cursorTime != nil && cursorID != nil {
		db = db.Where("(animes.aired_at < ? OR (animes.aired_at = ? AND animes.id < ?))", *cursorTime, *cursorTime, *cursorID)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}

	result := db.Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}

func (r GormAnimeCategoryUniverseRepository) GetByAnimeIdsAndCategoryUniverseIds(anime_ids []uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategoryUniverse, error) {
	var categoryAnime []entities.AnimeCategoryUniverse
	result := r.dbReplica.
		Where("anime_id IN (?) ", anime_ids).
		Where("category_universe_id IN (?)", category_ids).
		Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}
