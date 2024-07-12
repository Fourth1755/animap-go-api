package ports

import "github.com/Fourth1755/animap-go-api/internal/core/entities"

type UserRepository interface {
	Save(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetById(id uint) (*entities.User, error)
}
