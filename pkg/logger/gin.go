package logger

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// FromGin은 Gin 컨텍스트에서 로거를 가져옵니다
func FromGin(c *gin.Context) *slog.Logger {
	if l, exists := c.Get("logger"); exists {
		if logger, ok := l.(*slog.Logger); ok {
			return logger
		}
	}
	return Logger
}
