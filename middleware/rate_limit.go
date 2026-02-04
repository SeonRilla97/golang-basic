package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	// 전역 리미터 (모든 요청에 적용)
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "요청이 너무 많습니다. 잠시 후 다시 시도해주세요.",
				},
			})
			return
		}
		c.Next()
	}
}
