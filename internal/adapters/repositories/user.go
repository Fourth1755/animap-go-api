package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetById(id uuid.UUID) (*entities.User, error)
	UpdateUser(user *entities.User) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(user *entities.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormUserRepository) GetUserByEmail(email string) (*entities.User, error) {
	selectUser := new(entities.User)
	result := r.db.Where("email = ?", email).First(selectUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return selectUser, nil
}

func (r *GormUserRepository) GetById(id uuid.UUID) (*entities.User, error) {
	user := new(entities.User)
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *GormUserRepository) UpdateUser(user *entities.User) error {
	result := r.db.Model(&user).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
