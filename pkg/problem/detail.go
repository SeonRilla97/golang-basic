package problem

// Detail은 RFC 7807 Problem Details 형식입니다
type Detail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`

	// 확장 필드
	Code   string       `json:"code,omitempty"`   // 내부 에러 코드
	Errors []FieldError `json:"errors,omitempty"` // 필드별 에러
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Content-Type
const ContentType = "application/problem+json"
