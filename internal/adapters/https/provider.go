package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpProviderHandler struct {
	service services.ProviderService
}

func NewHttpProviderHandler(service services.ProviderService) *HttpProviderHandler {
	return &HttpProviderHandler{service: service}
}

func (h *HttpProviderHandler) CreateProvider(c *gin.Context) {
	var req dtos.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	provider, err := h.service.CreateProvider(req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, provider)
}

func (h *HttpProviderHandler) GetAllProviders(c *gin.Context) {
	providers, err := h.service.GetAllProviders()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, providers)
}

func (h *HttpProviderHandler) AddProviderToAnime(c *gin.Context) {
	var req dtos.AddProviderToAnimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.AddProviderToAnime(req); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Provider mapped to anime successfully"})
}
