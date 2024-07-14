package services

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
)

type AnimeCategoryService interface {
	AddAnimeToCategory(animeCategory *entities.AnimeCategory) error
}

type AnimeCategoryServiceImpl struct {
	repo ports.AnimeCategoryRepository
}

func NewAnimeCategoryService(repo ports.AnimeCategoryRepository) AnimeCategoryService {
	return &AnimeCategoryServiceImpl{repo: repo}
}

func (s AnimeCategoryServiceImpl) AddAnimeToCategory(animeCategory *entities.AnimeCategory) error {
	if err := s.repo.Save(animeCategory); err != nil {
		return err
	}
	return nil
}
