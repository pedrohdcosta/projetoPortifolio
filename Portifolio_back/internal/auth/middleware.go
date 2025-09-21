package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := []byte(os.Getenv("JWT_SECRET"))
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) { return secret, nil })
		if err != nil || !tok.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("sub", toString(claims["sub"]))
		c.Next()
	}
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}
