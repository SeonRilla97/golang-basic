package dto

import "time"

// CreateCommentRequest 댓글 생성 요청
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	Author   string `json:"author" binding:"required,max=50"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

// UpdateCommentRequest 댓글 수정 요청
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentResponse 댓글 응답
type CommentResponse struct {
	ID        uint               `json:"id"`
	PostID    uint               `json:"post_id"`
	ParentID  *uint              `json:"parent_id,omitempty"`
	Content   string             `json:"content"`
	Author    string             `json:"author"`
	CreatedAt time.Time          `json:"created_at"`
	Replies   []*CommentResponse `json:"replies,omitempty"`
}
