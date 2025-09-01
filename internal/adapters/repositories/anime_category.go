package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCategoryRepository interface {
	Save(animeCategory []entities.AnimeCategory) error
	GetByCategoryId(uuid.UUID) ([]GetByCategoryId, error)
	GetByAnimeIdAndCategoryIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategory, error)
}

type GormAnimeCategoryRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeCategoryRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeCategoryRepository {
	return &GormAnimeCategoryRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormAnimeCategoryRepository) Save(animeCategory []entities.AnimeCategory) error {
	if result := r.dbPrimary.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}

type GetByCategoryId struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Episodes int       `json:"episodes"`
	Seasonal string    `json:"seasonal"`
	Year     string    `json:"year"`
	Image    string    `json:"image"`
}

func (r GormAnimeCategoryRepository) GetByCategoryId(category_id uuid.UUID) ([]GetByCategoryId, error) {
	var animes []GetByCategoryId
	sql := `
		SELECT
			a.id,
			a.name AS name,
			a.episodes,
			a.seasonal,
			a.year,
			a.image AS image
		FROM
			animes a
		JOIN
			anime_categories ac ON a.id = ac.anime_id
		WHERE
			ac.category_id = ?
			AND a.is_show = true
		ORDER BY
			a.aired_at DESC`

	result := r.dbReplica.Raw(sql, category_id).Scan(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}

func (r GormAnimeCategoryRepository) GetByAnimeIdAndCategoryIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategory, error) {
	var categoryAnime []entities.AnimeCategory
	result := r.dbReplica.Where("anime_id = ?", anime_id).Where("category_id in (?)", category_ids).Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}
