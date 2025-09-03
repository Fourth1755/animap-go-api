package services

import (
	"context"
	"math"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/utils"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type CommentService interface {
	GetComments(animeID uuid.UUID, commentType string, page int, limit int) (*dtos.CommentAnimePaginatedResponse, error)
	CreateComment(ctx context.Context, req dtos.CreateCommentAnimeRequest) error
}

type commentServiceImpl struct {
	repo     repositories.CommentAnimeRepository
	userRepo repositories.UserRepository
}

func NewCommentService(repo repositories.CommentAnimeRepository, userRepo repositories.UserRepository) CommentService {
	return &commentServiceImpl{repo: repo, userRepo: userRepo}
}

func (s *commentServiceImpl) CreateComment(ctx context.Context, req dtos.CreateCommentAnimeRequest) error {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return errs.NewUnexpectedError()
	}
	userIdUuid := uuid.MustParse(userId)

	newComment := &entities.CommentAnime{
		ID:        uuid.New(),
		Message:   req.Message,
		Type:      "comment",
		AnimeID:   req.AnimeID,
		AuthorID:  userIdUuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.repo.Create(newComment)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *commentServiceImpl) GetComments(animeID uuid.UUID, commentType string, page int, limit int) (*dtos.CommentAnimePaginatedResponse, error) {
	// Get raw data from the repository
	queryResult, err := s.repo.GetByAnimeID(animeID, commentType, page, limit)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	// Map repository results to DTOs
	responseComments := make([]dtos.CommentAnimeResponse, len(queryResult.Results))
	for i, qr := range queryResult.Results {
		responseComments[i] = dtos.CommentAnimeResponse{
			ID:          qr.ID,
			Message:     qr.Message, // Corrected field mapping
			Type:        qr.Type,
			CreatedAt:   utils.TimeAgo(qr.CreatedAt),
			AuthorID:    qr.AuthorID,
			AuthorName:  qr.AuthorName,
			AuthorImage: qr.AuthorImage,
		}
	}

	// Construct the final paginated response
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(queryResult.TotalItems) / float64(limit)))
	}

	paginatedResponse := &dtos.CommentAnimePaginatedResponse{
		Data: responseComments,
		Pagination: dtos.PaginatedResponse{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: queryResult.TotalItems,
		},
	}

	return paginatedResponse, nil
}
