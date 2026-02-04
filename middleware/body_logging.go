package middleware

import (
	"bytes"
	"gorm-test/pkg/logger"
	"io"

	"github.com/gin-gonic/gin"
)

// bodyLogWriter는 응답 본문을 캡처합니다
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func BodyLogging(maxBodySize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.FromGin(c)

		// 요청 본문 읽기
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 본문을 다시 읽을 수 있도록 복원
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 요청 본문 로깅 (크기 제한)
		if len(requestBody) > 0 && len(requestBody) <= maxBodySize {
			log.Debug("요청 본문", "request_body", string(requestBody))
		}

		// 응답 캡처 준비
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = blw

		c.Next()

		// 응답 본문 로깅 (크기 제한)
		if blw.body.Len() > 0 && blw.body.Len() <= maxBodySize {
			log.Debug("응답 본문", "response_body", blw.body.String())
		}
	}
}
