package services

import (
	"os"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *entities.User) error
	Login(user *entities.User) (*dtos.LoginResponse, error)
}

type UserServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

const (
	TokenDuration = 24 * 3
)

func (s *UserServiceImpl) CreateUser(user *entities.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	uuid := uuid.NewV4()
	user.Password = string(hashPassword)
	user.UUID = uuid.String()
	err = s.repo.Save(user)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *UserServiceImpl) Login(user *entities.User) (*dtos.LoginResponse, error) {
	selectUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(user.Password))
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
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
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	loginResponse := dtos.LoginResponse{
		Token:  t,
		UserID: selectUser.UUID,
	}
	return &loginResponse, nil
}
