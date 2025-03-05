package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpStduioHandler struct {
	service services.StudioService
}

func NewHttpStduioHandler(service services.StudioService) *HttpStduioHandler {
	return &HttpStduioHandler{service: service}
}

func (h *HttpStduioHandler) GetAllStduio(c *gin.Context) {
	var studio dtos.StudioListRequest
	if err := c.BindJSON(&studio); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	response, err := h.service.GetAllStudio(studio)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
