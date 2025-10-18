package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCategoryUniverseRepository interface {
	Save(animeCategory []entities.AnimeCategoryUniverse) error
	GetByCategoryUniverseId(uuid.UUID) ([]entities.AnimeCategoryUniverse, error)
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

func (r GormAnimeCategoryUniverseRepository) GetByCategoryUniverseId(category_id uuid.UUID) ([]entities.AnimeCategoryUniverse, error) {
	var categoryAnime []entities.AnimeCategoryUniverse
	result := r.dbReplica.
		Joins("JOIN animes ON animes.id = anime_category_universes.anime_id").
		Preload("Anime").
		Where("anime_category_universes.category_universe_id = ?", category_id).
		Order("animes.aired_at DESC").
		Find(&categoryAnime)
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
