package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     sync.RWMutex
	r      rate.Limit
	b      int
	expiry time.Duration
}

func NewIPRateLimiter(r rate.Limit, b int, expiry time.Duration) *IPRateLimiter {
	irl := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		r:      r,
		b:      b,
		expiry: expiry,
	}

	// 오래된 항목 정리 (메모리 누수 방지)
	go irl.cleanupLoop()

	return irl
}

func (irl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	limiter, exists := irl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(irl.r, irl.b)
		irl.ips[ip] = limiter
	}

	return limiter
}

func (irl *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(irl.expiry)
	for range ticker.C {
		irl.mu.Lock()
		// 실제로는 마지막 접근 시간을 추적해야 함
		// 여기서는 간단히 전체 초기화
		irl.ips = make(map[string]*rate.Limiter)
		irl.mu.Unlock()
	}
}

func (irl *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := irl.getLimiter(ip)

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
