package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/pkg/jwt"
	"strings"
)

func Auth(j *jwt.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.String(http.StatusNotFound, "404 page not found")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.String(http.StatusNotFound, "404 page not found")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := j.ValidateToken(tokenStr)
		if err != nil {
			c.String(http.StatusNotFound, "404 page not found")
			c.Abort()
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
