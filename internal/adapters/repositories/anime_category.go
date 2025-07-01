package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCategoryRepository interface {
	Save(animeCategory []entities.AnimeCategory) error
	GetByCategoryId(uuid.UUID) ([]entities.AnimeCategory, error)
	GetByAnimeIdAndCategoryIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategory, error)
}

type GormAnimeCategoryRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeCategoryRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeCategoryRepository {
	return &GormAnimeCategoryRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormAnimeCategoryRepository) Save(animeCategory []entities.AnimeCategory) error {
	if result := r.dbPrimary.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormAnimeCategoryRepository) GetByCategoryId(category_id uuid.UUID) ([]entities.AnimeCategory, error) {
	var categoryAnime []entities.AnimeCategory
	result := r.dbReplica.Preload("Anime").Where("category_id = ?", category_id).Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}

func (r GormAnimeCategoryRepository) GetByAnimeIdAndCategoryIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategory, error) {
	var categoryAnime []entities.AnimeCategory
	result := r.dbReplica.Where("anime_id = ?", anime_id).Where("category_id in (?)", category_ids).Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}
