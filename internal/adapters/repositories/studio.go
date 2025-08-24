package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudioRepository interface {
	GetAll() ([]entities.Studio, error)
	GetByIds(ids []uuid.UUID) ([]entities.Studio, error)
	GetById(id uuid.UUID) (*entities.Studio, error)
	GetByMyAnimeListId(id int) (*entities.Studio, error)
	Save(studio entities.Studio) (*entities.Studio, error)
}

type GormStudioRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormStudioRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) StudioRepository {
	return &GormStudioRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormStudioRepository) GetAll() ([]entities.Studio, error) {
	var studio []entities.Studio
	if result := r.dbReplica.Find(&studio); result.Error != nil {
		return nil, result.Error
	}
	return studio, nil
}

func (r GormStudioRepository) GetByIds(ids []uuid.UUID) ([]entities.Studio, error) {
	var studio []entities.Studio
	if result := r.dbReplica.Where("id in (?)", ids).
		Find(&studio); result.Error != nil {
		return nil, result.Error
	}
	return studio, nil
}

func (r GormStudioRepository) GetById(id uuid.UUID) (*entities.Studio, error) {
	studio := new(entities.Studio)
	if result := r.dbReplica.First(&studio, id); result.Error != nil {
		return nil, result.Error
	}
	return studio, nil
}

func (r GormStudioRepository) GetByMyAnimeListId(id int) (*entities.Studio, error) {
	studio := new(entities.Studio)
	if result := r.dbReplica.Where("my_anime_list_id = ?", id).First(&studio); result.Error != nil {
		return nil, result.Error
	}
	return studio, nil
}

func (r GormStudioRepository) Save(studio entities.Studio) (*entities.Studio, error) {
	if result := r.dbPrimary.Create(&studio); result.Error != nil {
		return nil, result.Error
	}
	return &studio, nil
}
