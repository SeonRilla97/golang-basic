// Package logger는 slog 기반의 구조화된 로깅을 제공한다.
//
// Init으로 생성된 전역 로거(Logger)로부터 WithFields, WithContext를 통해
// 파생 로거를 만들 수 있다. 미들웨어에서 request_id 등을 포함한 파생 로거를
// context에 저장하면, 이후 핸들러·서비스에서는 FromContext만으로
// 별도 선언 없이 request_id가 자동으로 모든 로그에 기록된다.
//
//	// 미들웨어
//	reqLogger := logger.WithFields(map[string]any{"request_id": id})
//	ctx = logger.WithContext(ctx, reqLogger)
//
//	// 핸들러
//	logger.FromContext(ctx).Info("처리 완료")
//	// => level=INFO msg=처리완료 request_id=abc-123

// // 로그 샘플링의 경우 직접 구현 필요
// slog.Handler 인터페이스(4개 메서드)를 구현
// // 파일 입출력
// io.MultiWriter를 이용
package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

type Config struct {
	Level  string
	Pretty bool
}

func Init(cfg Config) {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	if cfg.Pretty {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Fatal(msg string, args ...any) {
	Logger.Error(msg, args...)
	os.Exit(1)
}
