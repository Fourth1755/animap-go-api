package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case errs.AppError:
		c.IndentedJSON(e.Code, gin.H{"message": e.Error()})
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
}
