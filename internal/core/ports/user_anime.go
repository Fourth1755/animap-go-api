package ports

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
)

type UserAnimeRepository interface {
	Save(userAnime *entities.UserAnime) error
	GetByUserId(id uint) ([]entities.UserAnime, error)
}
