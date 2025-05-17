package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCategoryUniverseRepository interface {
	Save(animeCategory []entities.AnimeCategoryUniverse) error
	GetByCategoryUniverseId(uuid.UUID) ([]entities.AnimeCategoryUniverse, error)
	GetByAnimeIdAndCategoryUniverseIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategoryUniverse, error)
}

type GormAnimeCategoryUniverseRepository struct {
	db *gorm.DB
}

func NewGormAnimeCategoryUniverseRepository(db *gorm.DB) AnimeCategoryUniverseRepository {
	return &GormAnimeCategoryUniverseRepository{db: db}
}

func (r GormAnimeCategoryUniverseRepository) Save(animeCategory []entities.AnimeCategoryUniverse) error {
	if result := r.db.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormAnimeCategoryUniverseRepository) GetByCategoryUniverseId(category_id uuid.UUID) ([]entities.AnimeCategoryUniverse, error) {
	var categoryAnime []entities.AnimeCategoryUniverse
	result := r.db.Preload("Anime").Where("category_id = ?", category_id).Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}

func (r GormAnimeCategoryUniverseRepository) GetByAnimeIdAndCategoryUniverseIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategoryUniverse, error) {
	var categoryAnime []entities.AnimeCategoryUniverse
	result := r.db.Where("anime_id = ?", anime_id).Where("category_universe_id in (?)", category_ids).Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}
