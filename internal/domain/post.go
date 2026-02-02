package domain

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Author    string         `gorm:"size:50;not null" json:"author"`
	Views     int            `gorm:"default:0" json:"views"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Post) TableName() string {
	return "posts"
}
