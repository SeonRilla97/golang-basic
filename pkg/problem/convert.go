package problem

import "gorm-test/pkg/apperror"

func FromAppError(appErr *apperror.AppError, instance string) *Detail {
	pd := &Detail{
		Type:     getTypeFromCode(appErr.Code),
		Title:    appErr.Message,
		Status:   appErr.HTTPStatus,
		Detail:   appErr.Detail,
		Code:     appErr.Code,
		Instance: instance,
	}

	// 필드 에러 변환
	if appErr.Fields != nil {
		pd.Errors = make([]FieldError, 0, len(appErr.Fields))
		for field, msg := range appErr.Fields {
			pd.Errors = append(pd.Errors, FieldError{
				Field:   field,
				Message: msg,
			})
		}
	}

	return pd
}

func getTypeFromCode(code string) string {
	switch code {
	case "NOT_FOUND":
		return TypeNotFound
	case "BAD_REQUEST":
		return TypeBadRequest
	case "VALIDATION_ERROR":
		return TypeValidation
	case "UNAUTHORIZED":
		return TypeUnauthorized
	case "FORBIDDEN":
		return TypeForbidden
	case "CONFLICT":
		return TypeConflict
	default:
		return TypeInternalError
	}
}
