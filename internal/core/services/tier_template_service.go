package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type TierTemplateService interface {
	GetAll() (*dtos.GetTierTemplatePaginatedResponse, error)
	Create(req dtos.CreateTierTemplateRequest) error
	GetById(id uuid.UUID) (*dtos.GetByIdTierTemplateResponse, error)
}

type tierTemplateService struct {
	repo repositories.TierTemplateRepository
}

func NewTierTemplateService(repo repositories.TierTemplateRepository) TierTemplateService {
	return &tierTemplateService{repo: repo}
}

func (s *tierTemplateService) GetAll() (*dtos.GetTierTemplatePaginatedResponse, error) {
	tierTemplates, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	var response []dtos.GetTierTemplateResponse
	for _, t := range tierTemplates {
		response = append(response, dtos.GetTierTemplateResponse{
			ID:          t.ID,
			Name:        t.Name,
			Type:        t.Type,
			PlayedCount: t.PlayedCount,
			TierList:    t.TierList,
			TotalItem:   t.TotalItem,
			IsPlay:      true,
			CreatedBy:   "Fourth",
		})
	}

	return &dtos.GetTierTemplatePaginatedResponse{
		Data: response,
	}, nil
}

func (s *tierTemplateService) Create(req dtos.CreateTierTemplateRequest) error {
	newUUID, err := uuid.NewV7()
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	tierTemplate := &entities.TierTemplate{
		ID:   newUUID,
		Name: req.Name,
		Type: req.Type,
	}

	if err := s.repo.Save(tierTemplate); err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *tierTemplateService) GetById(id uuid.UUID) (*dtos.GetByIdTierTemplateResponse, error) {
	tierTemplate, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Tier template not found")
	}
	itemList := []dtos.GetByIdTierTemplateResponseItem{}

	response := &dtos.GetByIdTierTemplateResponse{
		ID:          tierTemplate.ID,
		Name:        tierTemplate.Name,
		Type:        tierTemplate.Type,
		PlayedCount: tierTemplate.PlayedCount,
		TierList:    tierTemplate.TierList,
		TotalItem:   tierTemplate.TotalItem,
		Items:       itemList,
	}

	return response, nil
}
