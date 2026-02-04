package errors

import (
	"fmt"
	"runtime"
)

type AppError struct {
	Code    string         // 에러 코드
	Message string         // 사용자에게 보여줄 메시지
	Err     error          // 원본 에러
	Stack   []uintptr      // 스택 트레이스
	Context map[string]any // 추가 컨텍스트
}

func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Stack:   captureStack(),
		Context: make(map[string]any),
	}
}

func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   captureStack(),
		Context: make(map[string]any),
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithContext(key string, value any) *AppError {
	e.Context[key] = value
	return e
}

func captureStack() []uintptr {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	return pcs[:n]
}

func (e *AppError) StackTrace() []string {
	frames := runtime.CallersFrames(e.Stack)
	var trace []string

	for {
		frame, more := frames.Next()
		trace = append(trace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return trace
}
