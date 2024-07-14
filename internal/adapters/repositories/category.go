package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
	"gorm.io/gorm"
)

type GormCategoryRepository struct {
	db *gorm.DB
}

func NewGormCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &GormCategoryRepository{db: db}
}

func (r *GormCategoryRepository) Save(category *entities.Category) error {
	if result := r.db.Create(&category); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormCategoryRepository) GetAll() ([]entities.Category, error) {
	var categorise []entities.Category
	if result := r.db.Find(&categorise); result.Error != nil {
		return nil, result.Error
	}
	return categorise, nil
}

func (r *GormCategoryRepository) GetById(id uint) (*entities.Category, error) {
	category := new(entities.Category)
	if result := r.db.First(&category, id); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}
