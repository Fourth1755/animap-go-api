package ports

import "github.com/Fourth1755/animap-go-api/internal/core/entities"

type CategoryRepository interface {
	Save(category *entities.Category) error
	GetAll() ([]entities.Category, error)
	GetById(id uint) (*entities.Category, error)
}
