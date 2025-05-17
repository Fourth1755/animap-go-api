package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type CategoryUniverseService interface {
	CreateCategoryUniverse(category *entities.CategoryUniverse) error
	GetCategoryUniverses() ([]entities.CategoryUniverse, error)
	GetCategoryUniverseById(id uuid.UUID) (*entities.CategoryUniverse, error)
}

type CategoryUniverseServiceImpl struct {
	repo repositories.CategoryUniverseRepository
}

func NewCategoryUniverseService(repo repositories.CategoryUniverseRepository) CategoryUniverseService {
	return &CategoryUniverseServiceImpl{repo: repo}
}

func (s *CategoryUniverseServiceImpl) CreateCategoryUniverse(category *entities.CategoryUniverse) error {
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

func (s *CategoryUniverseServiceImpl) GetCategoryUniverses() ([]entities.CategoryUniverse, error) {
	category, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	return category, nil
}

func (s *CategoryUniverseServiceImpl) GetCategoryUniverseById(id uuid.UUID) (*entities.CategoryUniverse, error) {
	category, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("Category not found")
	}
	return category, nil
}
