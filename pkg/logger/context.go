package logger

import (
	"context"
	"log/slog"
)

// Package logger는 slog 기반의 구조화된 로깅을 제공한다.
//
// Init으로 생성된 전역 로거(Logger)로부터 WithFields, WithContext를 통해
// 파생 로거를 만들 수 있다. 미들웨어에서 request_id 등을 포함한 파생 로거를
// context에 저장하면, 이후 핸들러·서비스에서는 FromContext만으로
// 별도 선언 없이 request_id가 자동으로 모든 로그에 기록된다.
//
//	// 미들웨어
//	reqLogger := logger.WithFields(map[string]any{"request_id": id})
//	 ctx = logger.WithContext(ctx, reqLogger)
//
//	 // 핸들러
//	 logger.FromContext(ctx).Info("처리 완료")
//	 // => level=INFO msg=처리완료 request_id=abc-123
type ctxKey struct{}

// WithContext는 컨텍스트에 로거를 저장합니다.
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext는 컨텍스트에서 로거를 가져옵니다.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return Logger
}

// WithFields는 필드가 추가된 새 로거를 반환합니다.
func WithFields(fields map[string]any) *slog.Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return Logger.With(attrs...)
}
