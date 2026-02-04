package dto

import "strings"

// SortParams 정렬 파라미터
type SortParams struct {
	Sort string `form:"sort"` // 예: "created_at,desc" 또는 "views,desc|created_at,desc"
}

// SortItem 개별 정렬 조건
type SortItem struct {
	Field     string
	Direction string
}

// 허용된 정렬 필드
var allowedSortFields = map[string]bool{
	"id":         true,
	"title":      true,
	"author":     true,
	"views":      true,
	"created_at": true,
	"updated_at": true,
}

// Parse 정렬 문자열 파싱
func (s *SortParams) Parse() []SortItem {
	if s.Sort == "" {
		return []SortItem{{Field: "created_at", Direction: "DESC"}}
	}

	var items []SortItem

	// 다중 정렬: "|"로 구분
	parts := strings.Split(s.Sort, "|")
	for _, part := range parts {
		// 필드,방향 파싱
		tokens := strings.Split(strings.TrimSpace(part), ",")
		if len(tokens) < 1 {
			continue
		}

		field := strings.ToLower(strings.TrimSpace(tokens[0]))
		direction := "ASC"
		if len(tokens) > 1 {
			dir := strings.ToUpper(strings.TrimSpace(tokens[1]))
			if dir == "DESC" {
				direction = "DESC"
			}
		}

		// 허용된 필드만 추가
		if allowedSortFields[field] {
			items = append(items, SortItem{Field: field, Direction: direction})
		}
	}

	// 기본값
	if len(items) == 0 {
		return []SortItem{{Field: "created_at", Direction: "DESC"}}
	}

	return items
}

// ToOrderString GORM Order 문자열 생성
func (s *SortParams) ToOrderString() string {
	items := s.Parse()
	orders := make([]string, len(items))
	for i, item := range items {
		orders[i] = item.Field + " " + item.Direction
	}
	return strings.Join(orders, ", ")
}
