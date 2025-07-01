package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CharacterRepository interface {
	Save(character *entities.Character) (*entities.Character, error)
	GetById(id uuid.UUID) (*entities.Character, error)
	GetByIds(ids []uuid.UUID) ([]entities.Character, error)
}

type GormCharacterRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormCharacterRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) CharacterRepository {
	return &GormCharacterRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormCharacterRepository) Save(character *entities.Character) (*entities.Character, error) {
	if result := r.dbPrimary.Create(&character); result.Error != nil {
		return nil, result.Error
	}
	return character, nil
}

func (r *GormCharacterRepository) GetById(id uuid.UUID) (*entities.Character, error) {
	character := new(entities.Character)
	if result := r.dbReplica.First(&character, id); result.Error != nil {
		return nil, result.Error
	}
	return character, nil
}

func (r *GormCharacterRepository) GetByIds(ids []uuid.UUID) ([]entities.Character, error) {
	var characters []entities.Character
	if result := r.dbReplica.
		Where("id in (?)", ids).Find(&characters); result.Error != nil {
		return nil, result.Error
	}
	return characters, nil
}
