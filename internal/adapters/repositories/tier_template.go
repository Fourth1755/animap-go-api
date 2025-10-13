package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TierTemplateRepository interface {
	GetAll() ([]entities.TierTemplate, error)
	Save(*entities.TierTemplate) error
	GetById(id uuid.UUID) (*entities.TierTemplate, error)
}

type tierTemplateRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewTierTemplateRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) TierTemplateRepository {
	return &tierTemplateRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *tierTemplateRepository) GetAll() ([]entities.TierTemplate, error) {
	var tierTemplates []entities.TierTemplate
	if err := r.dbReplica.Find(&tierTemplates).Error; err != nil {
		return nil, err
	}
	return tierTemplates, nil
}

func (r *tierTemplateRepository) Save(tierTemplate *entities.TierTemplate) error {
	if result := r.dbPrimary.Create(&tierTemplate); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tierTemplateRepository) GetById(id uuid.UUID) (*entities.TierTemplate, error) {
	var tierTemplate entities.TierTemplate
	if err := r.dbReplica.First(&tierTemplate, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tierTemplate, nil
}
