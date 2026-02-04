package dto

import "time"

// CreatePostRequest 게시글 생성 요청
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=200,safe_string"`
	Content string `json:"content" binding:"required,safe_string"`
	URL     string `json:"url" binding:"omitempty,url,safe_url"`
}

// UpdatePostRequest 게시글 수정 요청
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
}

// PostResponse 게시글 응답
type PostResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostListResponse 게시글 목록 응답
type PostListResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	Highlight string    `json:"highlight,omitempty"` // 검색어 주변 텍스트 - FE 구현을 용이하게 하기 위함
}
