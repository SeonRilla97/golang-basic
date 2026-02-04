package safe

import (
	"log/slog"
	"runtime/debug"
)

// Go는 안전하게 고루틴을 실행합니다 -> Go Routine은 Recovery Middleware에서 잡지 못함 -> 별도 함수를 써서 고루틴을 사용하게 한다.
func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				slog.Error("고루틴 패닉 복구", "panic", err, "stack", stack)
			}
		}()
		fn()
	}()
}
