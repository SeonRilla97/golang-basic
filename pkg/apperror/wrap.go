package apperror

import "net/http"

// Wrap은 원본 에러를 래핑합니다
func Wrap(err error, code string, message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// WrapWithStatus는 상태 코드와 함께 래핑합니다
func WrapWithStatus(err error, status int, code string, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// WithDetail은 상세 설명을 추가합니다
func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

// WithFields는 필드 에러를 추가합니다
func (e *AppError) WithFields(fields map[string]string) *AppError {
	e.Fields = fields
	return e
}
