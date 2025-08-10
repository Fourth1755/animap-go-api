package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpUserHandler struct {
	service services.UserService
}

func NewHttpUserHandler(service services.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: service}
}

func (h *HttpUserHandler) CreateUser(c *gin.Context) {
	user := new(entities.User)
	if err := c.BindJSON(user); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.CreateUser(user)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Register success."})
}

func (h *HttpUserHandler) Login(c *gin.Context) {
	user := new(entities.User)
	if err := c.BindJSON(user); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	response, err := h.service.Login(user)
	if err != nil {
		handleError(c, errs.NewUnauthorizedError(err.Error()))
		return
	}
	c.SetCookie("jwt", response.Token, 3600*24, "/", "", false, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    response.Token,
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   false, // dev = false, prod = true
		SameSite: http.SameSiteLaxMode,
	})
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login success.",
		"token":   response.Token,
		"user_id": response.UserID,
	})
}

func (h *HttpUserHandler) GetUserInfo(c *gin.Context) {
	response, err := h.service.GetUserInfo(c)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *HttpUserHandler) GetUserByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	response, err := h.service.GetUserByUUID(uuid)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (h *HttpUserHandler) GetPresignedURLAvatar(c *gin.Context) {
	var req dtos.PresignUrlRequest
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	presignedURL, err := h.service.GetPresignedURLAvatar(c, &req)
	if err != nil {
		appError, ok := err.(errs.AppError)
		if ok {
			c.IndentedJSON(appError.Code, gin.H{"message": appError.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"url": presignedURL})
}

func (h *HttpUserHandler) UpdateUserInfo(c *gin.Context) {
	var request dtos.UpdateUserInfoRequest
	if err := c.BindJSON(&request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.UpdateUserInfo(c, &request); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update user info success."})
}

func (h *HttpUserHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // dev = false, prod = true
		SameSite: http.SameSiteLaxMode,
	})
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logout success."})
}
