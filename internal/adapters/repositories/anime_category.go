package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type AnimeCategoryRepository interface {
	Save(animeCategory *entities.AnimeCategory) error
}

type GormAnimeCategoryRepository struct {
	db *gorm.DB
}

func NewGormAnimeCategoryRepository(db *gorm.DB) AnimeCategoryRepository {
	return &GormAnimeCategoryRepository{db: db}
}

func (r GormAnimeCategoryRepository) Save(animeCategory *entities.AnimeCategory) error {
	if result := r.db.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}
