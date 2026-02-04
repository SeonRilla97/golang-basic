package middleware

import (
	"gorm-test/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// 모든 HTTP 요청을 로깅하는 미들웨어를 만들어 봅시다.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 요청 ID 가져오기 (RequestID 미들웨어가 먼저 실행되어야 함)
		requestID := c.GetString("requestID")

		// 요청별 로거 생성
		requestLogger := logger.Logger.With(
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"client_ip", c.ClientIP(),
		)

		// 컨텍스트에 로거 저장 (핸들러에서 사용 가능)
		c.Set("logger", requestLogger)

		// 요청 시작 로그
		requestLogger.Info("요청 시작",
			"user_agent", c.Request.UserAgent(),
			"query", query,
		)

		// 핸들러 실행
		c.Next()

		// 응답 완료 후 로그
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 로그 레벨 결정
		logFn := requestLogger.Info
		switch {
		case statusCode >= 500:
			logFn = requestLogger.Error
		case statusCode >= 400:
			logFn = requestLogger.Warn
		}

		logFn("요청 완료",
			"status", statusCode,
			"latency", latency,
			"body_size", c.Writer.Size(),
		)
	}
}
