package apperror

import (
	"fmt"
)

// AppError는 애플리케이션 에러를 나타냅니다
type AppError struct {
	HTTPStatus int               // HTTP 상태 코드
	Code       string            // 에러 코드 (예: USER_NOT_FOUND)
	Message    string            // 사용자에게 보여줄 메시지
	Detail     string            // 상세 설명 (개발자용)
	Err        error             // 원본 에러
	Fields     map[string]string // 필드별 에러 (유효성 검사)
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
