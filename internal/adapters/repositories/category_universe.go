package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryUniverseRepository interface {
	Save(category *entities.CategoryUniverse) error
	GetAll() ([]entities.CategoryUniverse, error)
	GetById(id uuid.UUID) (*entities.CategoryUniverse, error)
}

type GormCategoryUniverseRepository struct {
	db *gorm.DB
}

func NewGormCategoryUniverseRepository(db *gorm.DB) CategoryUniverseRepository {
	return &GormCategoryUniverseRepository{
		db: db,
	}
}

func (r *GormCategoryUniverseRepository) Save(category *entities.CategoryUniverse) error {
	if result := r.db.Create(&category); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormCategoryUniverseRepository) GetAll() ([]entities.CategoryUniverse, error) {
	var categorise []entities.CategoryUniverse
	if result := r.db.Find(&categorise); result.Error != nil {
		return nil, result.Error
	}
	return categorise, nil
}

func (r *GormCategoryUniverseRepository) GetById(id uuid.UUID) (*entities.CategoryUniverse, error) {
	category := new(entities.CategoryUniverse)
	if result := r.db.First(&category, id); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}
