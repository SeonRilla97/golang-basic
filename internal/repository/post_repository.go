package repository

import (
	"gorm-test/internal/domain"

	"gorm.io/gorm"
)

// PostRepository 게시글 저장소 인터페이스
type PostRepository interface {
	Create(post *domain.Post) error
	FindByID(id uint) (*domain.Post, error)
	FindAll(offset, limit int) ([]domain.Post, int64, error)
	Update(post *domain.Post) error
	Delete(id uint) error
	IncrementViews(id uint) error
}

// postRepository PostRepository 구현체
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository 생성자
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Create 게시글 생성
func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

// FindByID ID로 게시글 조회
func (r *postRepository) FindByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindAll 게시글 목록 조회 (페이징)
func (r *postRepository) FindAll(offset, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// 전체 개수 조회
	if err := r.db.Model(&domain.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 페이징 적용하여 조회
	err := r.db.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Update 게시글 수정
func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

// Delete 게시글 삭제
func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Post{}, id).Error
}

// IncrementViews 조회수 증가
func (r *postRepository) IncrementViews(id uint) error {
	return r.db.Model(&domain.Post{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).
		Error
}
