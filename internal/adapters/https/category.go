package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpCategoryHandler struct {
	service services.CategoryService
}

func NewHttpCategoryHandler(service services.CategoryService) *HttpCategoryHandler {
	return &HttpCategoryHandler{service: service}
}

func (h *HttpCategoryHandler) CreateCategory(c *gin.Context) {
	var category entities.Category
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
	c.JSON(http.StatusOK, category)
}

func (h *HttpCategoryHandler) GetCategoryById(c *gin.Context) {
	categoryId := c.Param("category_id")
	category, err := h.service.GetCategoryById(uuid.MustParse(categoryId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, category)
}
