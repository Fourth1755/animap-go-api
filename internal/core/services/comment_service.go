package services

import (
	"math"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type CommentService interface {
	GetComments(animeID uuid.UUID, commentType string, page int, limit int) (*dtos.CommentAnimePaginatedResponse, error)
	CreateComment(animeID uuid.UUID, authorID uuid.UUID, req dtos.CreateCommentAnimeRequest) error
}

type commentServiceImpl struct {
	repo     repositories.CommentAnimeRepository
	userRepo repositories.UserRepository
}

func NewCommentService(repo repositories.CommentAnimeRepository, userRepo repositories.UserRepository) CommentService {
	return &commentServiceImpl{repo: repo, userRepo: userRepo}
}

func (s *commentServiceImpl) CreateComment(animeID uuid.UUID, authorID uuid.UUID, req dtos.CreateCommentAnimeRequest) error {
	newComment := &entities.CommentAnime{
		ID:        uuid.New(),
		Message:   req.Message,
		Type:      "comment",
		AnimeID:   animeID,
		AuthorID:  authorID,
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
			ID:        qr.ID,
			Content:   qr.Message, // Corrected field mapping
			Type:      qr.Type,
			CreatedAt: qr.CreatedAt,
			Author: dtos.CommentAnimeAuthorResponse{
				ID:    qr.AuthorID,
				Name:  qr.AuthorName,
				Image: qr.AuthorImage,
			},
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
