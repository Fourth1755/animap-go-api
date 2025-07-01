package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *entities.User) error
	Login(user *entities.User) (*dtos.LoginResponse, error)
	GetUserInfo(ctx context.Context) (*dtos.GetUserInfoResponse, error)
	UpdateUserInfo(request *dtos.UpdateUserInfoRequest) error
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
	userId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	user.ID = userId
	user.Password = string(hashPassword)
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
		"uuid":    selectUser.ID,
		"picture": selectUser.ProfileImage,
		"name":    selectUser.Name,
		"email":   selectUser.Email,
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour * TokenDuration).Unix(),
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
		UserID: selectUser.ID,
	}
	return &loginResponse, nil
}

func (s *UserServiceImpl) GetUserInfo(ctx context.Context) (*dtos.GetUserInfoResponse, error) {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, errs.NewUnexpectedError()
	}
	fmt.Println(userId)
	userIdUuid := uuid.MustParse(userId)
	user, err := s.repo.GetById(userIdUuid)
	if err != nil {
		return nil, err
	}
	return &dtos.GetUserInfoResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		Description:  user.Description,
	}, nil
}

func (s *UserServiceImpl) UpdateUserInfo(request *dtos.UpdateUserInfoRequest) error {
	_, err := s.repo.GetById(request.ID)
	if err != nil {
		logs.Error(err)
		return errs.NewNotFoundError("User not found")
	}

	user := entities.User{
		Name:         request.Name,
		Email:        request.Email,
		ProfileImage: request.ProfileImage,
		Description:  request.Description,
	}

	if err := s.repo.UpdateUser(&user); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}
