package apperror

import (
	"errors"
	"net/http"
)

// IsAppError는 에러가 AppError인지 확인합니다
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError는 에러를 AppError로 변환합니다
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetHTTPStatus는 에러에서 HTTP 상태 코드를 추출합니다
func GetHTTPStatus(err error) int {
	if appErr, ok := AsAppError(err); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
