package problem

import "net/http"

func NotFound(resource string, instance string) *Detail {
	return &Detail{
		Type:     TypeNotFound,
		Title:    resource + "을(를) 찾을 수 없습니다",
		Status:   http.StatusNotFound,
		Code:     "NOT_FOUND",
		Instance: instance,
	}
}

func BadRequest(message string, instance string) *Detail {
	return &Detail{
		Type:     TypeBadRequest,
		Title:    "잘못된 요청입니다",
		Status:   http.StatusBadRequest,
		Detail:   message,
		Code:     "BAD_REQUEST",
		Instance: instance,
	}
}

func Validation(errors []FieldError, instance string) *Detail {
	return &Detail{
		Type:     TypeValidation,
		Title:    "입력값이 올바르지 않습니다",
		Status:   http.StatusBadRequest,
		Detail:   "요청 본문의 유효성 검사에 실패했습니다",
		Code:     "VALIDATION_ERROR",
		Instance: instance,
		Errors:   errors,
	}
}

func Unauthorized(message string, instance string) *Detail {
	if message == "" {
		message = "인증이 필요합니다"
	}
	return &Detail{
		Type:     TypeUnauthorized,
		Title:    message,
		Status:   http.StatusUnauthorized,
		Code:     "UNAUTHORIZED",
		Instance: instance,
	}
}

func Forbidden(message string, instance string) *Detail {
	if message == "" {
		message = "접근 권한이 없습니다"
	}
	return &Detail{
		Type:     TypeForbidden,
		Title:    message,
		Status:   http.StatusForbidden,
		Code:     "FORBIDDEN",
		Instance: instance,
	}
}

func InternalError(instance string) *Detail {
	return &Detail{
		Type:     TypeInternalError,
		Title:    "서버 내부 오류가 발생했습니다",
		Status:   http.StatusInternalServerError,
		Detail:   "잠시 후 다시 시도해 주세요",
		Code:     "INTERNAL_ERROR",
		Instance: instance,
	}
}
