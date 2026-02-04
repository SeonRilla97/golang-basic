package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
	// 안전한 문자열 (스크립트 태그 없음)
	v.RegisterValidation("safe_string", validateSafeString)

	// 알파벳과 숫자만
	v.RegisterValidation("alphanum_unicode", validateAlphanumUnicode)

	// 안전한 URL
	v.RegisterValidation("safe_url", validateSafeURL)
}

func validateSafeString(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	// 스크립트 태그 검사
	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	if scriptPattern.MatchString(s) {
		return false
	}

	// 이벤트 핸들러 검사
	eventPattern := regexp.MustCompile(`(?i)on\w+\s*=`)
	if eventPattern.MatchString(s) {
		return false
	}

	return true
}

func validateAlphanumUnicode(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' {
			return false
		}
	}
	return true
}

func validateSafeURL(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	// javascript: 프로토콜 금지
	if regexp.MustCompile(`(?i)^javascript:`).MatchString(s) {
		return false
	}

	// data: 프로토콜 금지
	if regexp.MustCompile(`(?i)^data:`).MatchString(s) {
		return false
	}

	return true
}
