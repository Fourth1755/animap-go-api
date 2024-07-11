package services

import (
	"os"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *entities.User) error
	Login(user *entities.User) (string, error)
}

type UserServiceImpl struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

const (
	TokenDuration = 24 * 3
)

func (s *UserServiceImpl) CreateUser(user *entities.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashPassword)
	err = s.repo.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) Login(user *entities.User) (string, error) {
	selectUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"email": selectUser.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * TokenDuration).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}
