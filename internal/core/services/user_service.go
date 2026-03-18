package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/aws"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserService interface {
	CreateUser(user *entities.User) error
	Login(user *entities.User) (*dtos.LoginResponse, error)
	LoginWithGoogle(code string) (*dtos.LoginResponse, error)
	GetGoogleOAuthURL(state string) string
	GetGoogleFrontendURL() string
	GetUserInfo(ctx context.Context) (*dtos.GetUserInfoResponse, error)
	UpdateUserInfo(ctx context.Context, request *dtos.UpdateUserInfoRequest) error
	GetUserByUUID(uuid string) (*dtos.GetUserInfoResponse, error)
	GetPresignedURLAvatar(ctx context.Context, req *dtos.PresignUrlRequest) (string, error)
}

type UserServiceImpl struct {
	repo          repositories.UserRepository
	s3Service     aws.S3Service
	configService config.ConfigService
}

func NewUserService(
	repo repositories.UserRepository,
	s3Service aws.S3Service,
	configService config.ConfigService) UserService {
	return &UserServiceImpl{
		repo:          repo,
		s3Service:     s3Service,
		configService: configService}
}

const (
	TokenDuration = 24 * 3
)

func (s *UserServiceImpl) GetPresignedURLAvatar(ctx context.Context, req *dtos.PresignUrlRequest) (string, error) {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return "", errs.NewUnexpectedError()
	}
	fmt.Println("userId " + userId)
	key := fmt.Sprintf("user/user-avatar/%s/%s", userId, req.FileName)
	presignedURL, err := s.s3Service.GetPresignedURL(s.configService.GetAWS().S3Bucket, key)
	if err != nil {
		logs.Error(err)
		return "", errs.NewUnexpectedError()
	}
	return presignedURL, nil
}

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
		"role":    "user",
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

func (s *UserServiceImpl) UpdateUserInfo(ctx context.Context, request *dtos.UpdateUserInfoRequest) error {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return errs.NewUnexpectedError()
	}

	userIdUuid := uuid.MustParse(userId)
	_, err := s.repo.GetById(userIdUuid)
	if err != nil {
		logs.Error(err)
		return errs.NewNotFoundError("User not found")
	}

	user := entities.User{
		ID:           userIdUuid,
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

func (s *UserServiceImpl) GetUserByUUID(uuidStr string) (*dtos.GetUserInfoResponse, error) {
	uuidParsed := uuid.MustParse(uuidStr)
	user, err := s.repo.GetById(uuidParsed)
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

func (s *UserServiceImpl) googleOAuthConfig() *oauth2.Config {
	cfg := s.configService.GetGoogleOAuth()
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (s *UserServiceImpl) GetGoogleOAuthURL(state string) string {
	return s.googleOAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *UserServiceImpl) GetGoogleFrontendURL() string {
	return s.configService.GetGoogleOAuth().FrontendURL
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (s *UserServiceImpl) LoginWithGoogle(code string) (*dtos.LoginResponse, error) {
	oauthCfg := s.googleOAuthConfig()

	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnauthorizedError("failed to exchange google oauth code")
	}

	client := oauthCfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	var googleUser googleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	// Find existing user by GoogleID, then fallback to email
	user, err := s.repo.GetUserByGoogleID(googleUser.ID)
	if err != nil {
		// Not found by GoogleID — try by email
		user, err = s.repo.GetUserByEmail(googleUser.Email)
		if err != nil {
			// Create new user
			newID, idErr := uuid.NewV7()
			if idErr != nil {
				logs.Error(idErr)
				return nil, errs.NewUnexpectedError()
			}
			user = &entities.User{
				ID:           newID,
				Email:        googleUser.Email,
				Name:         googleUser.Name,
				GoogleID:     googleUser.ID,
				ProfileImage: googleUser.Picture,
			}
			if saveErr := s.repo.Save(user); saveErr != nil {
				logs.Error(saveErr)
				return nil, errs.NewUnexpectedError()
			}
		} else {
			// Link GoogleID to existing email account
			user.GoogleID = googleUser.ID
			if updateErr := s.repo.UpdateUser(user); updateErr != nil {
				logs.Error(updateErr)
			}
		}
	}

	// Issue JWT
	claims := jwt.MapClaims{
		"uuid":    user.ID,
		"picture": user.ProfileImage,
		"name":    user.Name,
		"email":   user.Email,
		"role":    "user",
		"exp":     time.Now().Add(time.Hour * TokenDuration).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	return &dtos.LoginResponse{Token: t, UserID: user.ID}, nil
}
