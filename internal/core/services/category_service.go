package services

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
)

type CategoryService interface {
	CreateCategory(category *entities.Category) error
	Getcategorise() ([]entities.Category, error)
	GetCategoryById(id uint) (*entities.Category, error)
}

type CategoryServiceImpl struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) CreateCategory(category *entities.Category) error {
	if err := s.repo.Save(category); err != nil {
		return err
	}
	return nil
}

func (s *CategoryServiceImpl) Getcategorise() ([]entities.Category, error) {
	category, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryServiceImpl) GetCategoryById(id uint) (*entities.Category, error) {
	category, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return category, nil
}
