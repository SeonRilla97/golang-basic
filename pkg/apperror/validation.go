package apperror

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FromValidationErrors는 validator 에러를 AppError로 변환합니다
func FromValidationErrors(err error) *AppError {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		fields := make(map[string]string)
		for _, e := range validationErrs {
			fields[e.Field()] = translateValidationError(e)
		}
		return ValidationError(fields)
	}
	return BadRequest("입력값이 올바르지 않습니다")
}

func translateValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "필수 항목입니다"
	case "email":
		return "올바른 이메일 형식이 아닙니다"
	case "min":
		return fmt.Sprintf("최소 %s자 이상이어야 합니다", e.Param())
	case "max":
		return fmt.Sprintf("최대 %s자까지 입력 가능합니다", e.Param())
	case "gte":
		return fmt.Sprintf("%s 이상이어야 합니다", e.Param())
	case "lte":
		return fmt.Sprintf("%s 이하여야 합니다", e.Param())
	default:
		return "올바르지 않은 값입니다"
	}
}
