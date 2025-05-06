package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/gin-gonic/gin"
)

type HttpCommonHandler struct {
	service services.CommonService
}

func NewHttpCommonHandler(service services.CommonService) *HttpCommonHandler {
	return &HttpCommonHandler{service: service}
}

func (h *HttpCommonHandler) GetSeasonalAndYear(c *gin.Context) {
	response, err := h.service.GetSeasonalAndYear()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
