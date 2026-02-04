package logger

import (
	"gorm-test/pkg/apperror"
	"log/slog"
)

// LogError는 에러를 구조화된 형태로 로깅합니다
func LogError(log *slog.Logger, err error) {
	// AppError인 경우 추가 정보 포함
	if appErr, ok := err.(*apperror.apperror); ok {
		attrs := []any{
			"error_code", appErr.Code,
			"context", appErr.Context,
			"stack_trace", appErr.StackTrace(),
		}
		if appErr.Err != nil {
			attrs = append(attrs, "error", appErr.Err)
		}
		log.Error("에러 발생", attrs...)
	} else {
		log.Error("에러 발생", "error", err)
	}
}
