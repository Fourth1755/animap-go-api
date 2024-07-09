package ports

import (
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
)

type AnimeRepository interface {
	Save(anime entities.Anime) error
	GetById(id int) (*entities.Anime, error)
	GetAll(query dtos.AnimeQueryDTO) ([]entities.Anime, error)
}
