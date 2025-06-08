package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCharacterRepository interface {
	Save(category *entities.AnimeCharacter) error
	GetByAnimeId(animeId uuid.UUID) ([]entities.AnimeCharacter, error)
}

type GormAnimeCharacterRepository struct {
	db *gorm.DB
}

func NewGormAnimeCharacterRepository(db *gorm.DB) AnimeCharacterRepository {
	return &GormAnimeCharacterRepository{db: db}
}

func (r *GormAnimeCharacterRepository) Save(character *entities.AnimeCharacter) error {
	if result := r.db.Create(&character); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeCharacterRepository) GetByAnimeId(animeId uuid.UUID) ([]entities.AnimeCharacter, error) {
	var animeCharacter []entities.AnimeCharacter
	result := r.db.Preload("Character").Where("anime_id = ?", animeId).Find(&animeCharacter)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeCharacter, nil
}
