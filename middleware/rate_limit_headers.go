package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 클라이언트에게 남은 요청 수를 알려줍니다.
type RateLimitConfig struct {
	Rate   rate.Limit    // 초당 요청 수
	Burst  int           // 버스트 크기
	Window time.Duration // 시간 창
}

func RateLimiterWithHeaders(cfg RateLimitConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(cfg.Rate, cfg.Burst)

	return func(c *gin.Context) {
		// Rate Limit 헤더 설정
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.Burst))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%.0f", limiter.Tokens()))

		if !limiter.Allow() {
			reservation := limiter.Reserve()
			delay := reservation.Delay()
			reservation.Cancel()

			c.Header("Retry-After", strconv.Itoa(int(delay.Seconds())+1))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "RATE_LIMITED",
					"message":     "요청이 너무 많습니다.",
					"retry_after": int(delay.Seconds()) + 1,
				},
			})
			return
		}

		c.Next()
	}
}
