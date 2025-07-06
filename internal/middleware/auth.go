package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthRequired(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Login required"})
		return
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Token"})
		return
	}
	claim := token.Claims.(jwt.MapClaims)
	c.Set("userId", claim["uuid"])
	c.Next()
}
