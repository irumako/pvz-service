package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pvz-service/internal/entity"
	"pvz-service/pkg/jwt"
)

func Role(allowedRoles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsValue, exists := c.Get("userClaims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "User claims missing from context",
			})
			return
		}

		claims, ok := claimsValue.(*jwt.CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Invalid token claims",
			})
			return
		}

		allowed := false
		for _, role := range allowedRoles {
			if string(role) == claims.Role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Доступ запрещен",
			})
			return
		}

		c.Next()
	}
}
