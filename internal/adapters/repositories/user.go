package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
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

func (r *GormUserRepository) GetById(id uint) (*entities.User, error) {
	user := new(entities.User)
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
