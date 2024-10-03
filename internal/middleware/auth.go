package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthRequired(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.String(http.StatusNotFound, "Cookie not found")
		return
	}
	jwtSecretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, "")
		return
	}
	claim := token.Claims.(jwt.MapClaims)
	fmt.Println(claim["role"])
	c.Next()
}
