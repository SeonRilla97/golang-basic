package dto

// Response API 응답 래퍼
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo 에러 정보
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta 페이징 메타 정보
type Meta struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SuccessResponse 성공 응답 생성
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// SuccessWithMeta 페이징 포함 성공 응답 생성
func SuccessWithMeta(data interface{}, meta *Meta) Response {
	return Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// ErrorResponse 에러 응답 생성
func ErrorResponse(code, message string) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}
