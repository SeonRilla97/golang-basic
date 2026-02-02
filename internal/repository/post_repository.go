package repository

import (
	"gorm-test/internal/domain"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func (r *PostRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}
