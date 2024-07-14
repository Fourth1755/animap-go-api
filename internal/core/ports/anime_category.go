package ports

import "github.com/Fourth1755/animap-go-api/internal/core/entities"

type AnimeCategoryRepository interface {
	Save(animeCategory *entities.AnimeCategory) error
}
