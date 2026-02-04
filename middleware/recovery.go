package middleware

import (
	"errors"
	"gorm-test/pkg/logger"
	"gorm-test/pkg/problem"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"syscall"

	"github.com/gin-gonic/gin"
)

type RecoveryConfig struct {
	EnableStackTrace bool                        // 스택 트레이스 로깅 여부
	NotifyFunc       func(err any, stack string) // 외부 알림 함수
}

func isBrokenPipe(err any) bool {
	if netErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			if errors.Is(sysErr.Err, syscall.EPIPE) ||
				errors.Is(sysErr.Err, syscall.ECONNRESET) {
				return true
			}
		}
	}
	return false
}

func Recovery(cfg RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.FromGin(c)
				// 스택 트레이스 캡처
				stack := string(debug.Stack())
				if isBrokenPipe(err) {
					log.Warn("브로큰 파이프", "error", err)
					c.Abort()
					return
				}

				// 로깅

				attrs := []any{
					"panic", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				}
				if cfg.EnableStackTrace {
					attrs = append(attrs, "stack", stack)
				}
				log.Error("패닉 복구", attrs...)

				// 외부 알림 (Slack, Discord, Sentry 등)
				if cfg.NotifyFunc != nil {
					go cfg.NotifyFunc(err, stack)
				}

				// RFC 7807 형식 응답
				pd := problem.InternalError(c.Request.URL.Path)
				c.Header("Content-Type", problem.ContentType)
				c.AbortWithStatusJSON(http.StatusInternalServerError, pd)
			}
		}()

		c.Next()
	}
}

// DefaultRecovery는 기본 설정의 Recovery 미들웨어입니다
func DefaultRecovery() gin.HandlerFunc {
	return Recovery(RecoveryConfig{
		EnableStackTrace: true,
	})
}
