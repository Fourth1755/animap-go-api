package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Save(category entities.Category) (*entities.Category, error)
	GetAll() ([]entities.Category, error)
	GetById(id uuid.UUID) (*entities.Category, error)
	GetByMyAnimeListId(id int) (*entities.Category, error)
}

type GormCategoryRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormCategoryRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) CategoryRepository {
	return &GormCategoryRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormCategoryRepository) Save(category entities.Category) (*entities.Category, error) {
	if result := r.dbPrimary.Create(&category); result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (r *GormCategoryRepository) GetAll() ([]entities.Category, error) {
	var categorise []entities.Category
	if result := r.dbReplica.Find(&categorise); result.Error != nil {
		return nil, result.Error
	}
	return categorise, nil
}

func (r *GormCategoryRepository) GetById(id uuid.UUID) (*entities.Category, error) {
	category := new(entities.Category)
	if result := r.dbReplica.First(&category, id); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (r *GormCategoryRepository) GetByMyAnimeListId(id int) (*entities.Category, error) {
	category := new(entities.Category)
	if result := r.dbReplica.Where("my_anime_list_id = ?", id).First(&category); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}
