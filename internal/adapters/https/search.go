package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpSearchHandler struct {
	service services.SearchService
}

func NewHttpSearchHandler(service services.SearchService) *HttpSearchHandler {
	return &HttpSearchHandler{service: service}
}

func (h *HttpSearchHandler) Search(c *gin.Context) {
	var req dtos.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	result, err := h.service.Search(req.Keyword)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
