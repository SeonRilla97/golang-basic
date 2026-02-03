package apperror

import (
	"fmt"
	"net/http"
)

// NotFound는 리소스를 찾을 수 없을 때 사용합니다
func NotFound(resource string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s을(를) 찾을 수 없습니다", resource),
	}
}

// NotFoundWithID는 ID와 함께 Not Found 에러를 생성합니다
func NotFoundWithID(resource string, id any) *AppError {
	return &AppError{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s을(를) 찾을 수 없습니다", resource),
		Detail:     fmt.Sprintf("ID: %v", id),
	}
}

// BadRequest는 잘못된 요청일 때 사용합니다
func BadRequest(message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

// ValidationError는 유효성 검사 실패 시 사용합니다
func ValidationError(fields map[string]string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    "입력값이 올바르지 않습니다",
		Fields:     fields,
	}
}

// Unauthorized는 인증이 필요할 때 사용합니다
func Unauthorized(message string) *AppError {
	if message == "" {
		message = "인증이 필요합니다"
	}
	return &AppError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

// Forbidden은 권한이 없을 때 사용합니다
func Forbidden(message string) *AppError {
	if message == "" {
		message = "접근 권한이 없습니다"
	}
	return &AppError{
		HTTPStatus: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

// Conflict는 리소스 충돌 시 사용합니다
func Conflict(message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

// InternalError는 서버 내부 오류일 때 사용합니다
func InternalError(err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "서버 내부 오류가 발생했습니다",
		Err:        err,
	}
}
