package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole은 특정 역할만 허용하는 미들웨어입니다.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool)
	for _, role := range allowedRoles {
		roleSet[role] = true
	}

	return func(c *gin.Context) {
		claims, ok := GetCurrentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "인증이 필요합니다",
			})
			return
		}

		if !roleSet[claims.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "권한이 없습니다",
			})
			return
		}

		c.Next()
	}
}
