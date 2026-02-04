package middleware

import (
	"bytes"
	"encoding/json"
	"gorm-test/pkg/sanitize"
	"io"

	"github.com/gin-gonic/gin"
)

func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JSON 요청만 처리
		if c.ContentType() != "application/json" {
			c.Next()
			return
		}

		// 본문 읽기
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		// JSON 파싱
		var data map[string]any
		if err := json.Unmarshal(body, &data); err != nil {
			// 파싱 실패 시 원본 유지
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}

		// 문자열 필드 정제
		sanitizeMap(data)

		// 정제된 본문으로 교체
		sanitized, _ := json.Marshal(data)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitized))

		c.Next()
	}
}

func sanitizeMap(m map[string]any) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			m[k] = sanitize.Strict(val)
		case map[string]any:
			sanitizeMap(val)
		case []any:
			sanitizeSlice(val)
		}
	}
}

func sanitizeSlice(s []any) {
	for i, v := range s {
		switch val := v.(type) {
		case string:
			s[i] = sanitize.Strict(val)
		case map[string]any:
			sanitizeMap(val)
		}
	}
}
