package adapters

import (
	"net/http"
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpCategoryHandler struct {
	service services.CategoryService
}

func NewHttpCategoryHandler(service services.CategoryService) *HttpCategoryHandler {
	return &HttpCategoryHandler{service: service}
}

func (h *HttpCategoryHandler) CreateCategory(c *gin.Context) {
	category := new(entities.Category)
	if err := c.BindJSON(&category); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.CreateCategory(category); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create category success"})
}

func (h *HttpCategoryHandler) Getcategorise(c *gin.Context) {
	category, err := h.service.Getcategorise()
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *HttpCategoryHandler) GetCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	category, err := h.service.GetCategoryById(uint(id))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, category)
}
