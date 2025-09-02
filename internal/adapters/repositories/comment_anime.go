package repositories

import (
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QueryResult holds the raw, flat result from the database query
type CommentQueryResult struct {
	ID          uuid.UUID
	Message     string
	Type        string
	CreatedAt   time.Time
	AuthorID    uuid.UUID
	AuthorName  string
	AuthorImage *string
}

// PaginatedCommentQueryResult is the repository's return type
// It contains the raw data and total count for the service layer to process.
type PaginatedCommentQueryResult struct {
	Results    []CommentQueryResult
	TotalItems int64
}

type CommentAnimeRepository interface {
	GetByAnimeID(animeID uuid.UUID, commentType string, page int, limit int) (*PaginatedCommentQueryResult, error)
}

type GormCommentAnimeRepository struct {
	dbReplica *gorm.DB
}

func NewGormCommentAnimeRepository(dbReplica *gorm.DB) CommentAnimeRepository {
	return &GormCommentAnimeRepository{dbReplica: dbReplica}
}

func (r *GormCommentAnimeRepository) GetByAnimeID(animeID uuid.UUID, commentType string, page int, limit int) (*PaginatedCommentQueryResult, error) {
	var queryResults []CommentQueryResult
	var totalItems int64

	offset := (page - 1) * limit

	// Query for the total count
	countQuery := r.dbReplica.Model(&entities.CommentAnime{}).Where("anime_id = ?", animeID)
	if commentType != "" {
		countQuery = countQuery.Where("type = ?", commentType)
	}
	if err := countQuery.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	// Query for the paginated data
	dataQuery := r.dbReplica.Model(&entities.CommentAnime{}).
		Select("comments.id, comments.message, comments.type, comments.created_at, users.id as author_id, users.name as author_name, users.profile_image as author_image").
		Joins("join users on users.id = comments.author_id").
		Where("comments.anime_id = ?", animeID)

	if commentType != "" {
		dataQuery = dataQuery.Where("comments.type = ?", commentType)
	}

	result := dataQuery.Order("comments.created_at desc").
		Offset(offset).
		Limit(limit).
		Scan(&queryResults)

	if result.Error != nil {
		return nil, result.Error
	}

	return &PaginatedCommentQueryResult{
		Results:    queryResults,
		TotalItems: totalItems,
	}, nil
}
