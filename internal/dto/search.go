package dto

// SearchParams 검색 파라미터
type SearchParams struct {
	Query      string `form:"q"`                                                // 검색어
	SearchType string `form:"type" binding:"omitempty,oneof=title content all"` // 검색 유형
}

// 검색 유형 상수
const (
	SearchTypeTitle   = "title"
	SearchTypeContent = "content"
	SearchTypeAll     = "all"
)

// GetSearchType 검색 유형 반환 (기본값: all)
func (s *SearchParams) GetSearchType() string {
	if s.SearchType == "" {
		return SearchTypeAll
	}
	return s.SearchType
}
