package dto

import "time"

// CreatePostRequest 게시글 생성 요청
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required,max=50"`
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
}
