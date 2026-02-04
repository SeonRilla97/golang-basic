package domain

import (
	"time"

	"gorm.io/gorm"
)

// Comment 댓글 도메인 모델
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PostID    uint           `gorm:"not null;index" json:"post_id"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"` // 최상위 댓글의 경우 nil로 부모 없음을 표현한다.
	Content   string         `gorm:"type:text;not null" json:"content"`
	Author    string         `gorm:"size:50;not null" json:"author"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 연관관계
	Post    Post      `gorm:"foreignKey:PostID" json:"-"`                   //다대일 (응답에서 제외)
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"-"`                 // 자기참조 (응답에서 제외)
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"` // 일대다
}

// TableName 테이블 이름 지정
func (Comment) TableName() string {
	return "comments"
}
