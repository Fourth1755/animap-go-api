package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TierTemplateHandler struct {
	service services.TierTemplateService
}

func NewTierTemplateHandler(service services.TierTemplateService) *TierTemplateHandler {
	return &TierTemplateHandler{service: service}
}

func (h *TierTemplateHandler) GetAll(c *gin.Context) {
	tierTemplates, err := h.service.GetAll()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tierTemplates)
}

func (h *TierTemplateHandler) Create(c *gin.Context) {
	var req dtos.CreateTierTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.Create(req); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tier template created successfully"})
}

func (h *TierTemplateHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	tierTemplate, err := h.service.GetById(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tierTemplate)
}
