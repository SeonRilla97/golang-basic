package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 요청 별 고유 ID 부착

var RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 헤더에서 요청 ID 확인 (외부에서 전달된 경우)
		requestID := c.GetHeader(RequestIDHeader)

		// 없으면 새로 생성
		if requestID == "" {
			requestID = uuid.New().String()

		}

		// 컨텍스트와 응답 헤더에 설정
		c.Set("requestID", requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}
