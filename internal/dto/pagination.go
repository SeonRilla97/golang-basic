package dto

/**
Pagination 페이징 요청

단점
- OFFSET 값이 커지면 DB는 그만큼 데이터를 건너뛰어서 가져온다
- 페이징 도중 새 글 추가 시 두 페이지에 노출되거나, 누락될 수 있음

=> 커서기반 필요 이유
*/

type Pagination struct {
	Page int `form:"page" binding:"min=1"`
	Size int `form:"size" binding:"min=1,max=100"`
}

// Offset 오프셋 계산
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

// NewPagination 기본값 적용
func NewPagination(page, size, defaultSize, maxSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}
	return &Pagination{Page: page, Size: size}
}
