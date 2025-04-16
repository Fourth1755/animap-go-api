package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(category *entities.Category) error
	Getcategorise() ([]entities.Category, error)
	GetCategoryById(id uint) (*entities.Category, error)
}

type CategoryServiceImpl struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) CreateCategory(category *entities.Category) error {
	categoryId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	category.ID = categoryId
	if err := s.repo.Save(category); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *CategoryServiceImpl) Getcategorise() ([]entities.Category, error) {
	category, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	return category, nil
}

func (s *CategoryServiceImpl) GetCategoryById(id uint) (*entities.Category, error) {
	category, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("Category not found")
	}
	return category, nil
}
