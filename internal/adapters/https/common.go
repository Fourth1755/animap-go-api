package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpCommonHandler struct {
	service services.CommonService
}

func NewHttpCommonHandler(service services.CommonService) *HttpCommonHandler {
	return &HttpCommonHandler{service: service}
}

func (h *HttpCommonHandler) GetSeasonalAndYear(c *gin.Context) {
	request := new(dtos.GetSeasonalAndYearRequest)
	if err := c.BindJSON(request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	response, err := h.service.GetSeasonalAndYear(*request)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
