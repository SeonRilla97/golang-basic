package repository

import (
	"gorm-test/internal/domain"

	"gorm.io/gorm"
)

// CommentRepository 댓글 저장소 인터페이스
type CommentRepository interface {
	Create(comment *domain.Comment) error
	FindByID(id uint) (*domain.Comment, error)
	FindByPostID(postID uint) ([]domain.Comment, error)
	FindByPostIDWithReplies(postID uint) ([]domain.Comment, error)
	Update(comment *domain.Comment) error
	Delete(id uint) error
	HasReplies(commentID uint) (bool, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByPostID(postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID). // 최상위 댓글만
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByPostIDWithReplies 게시글의 댓글 목록 조회 (대댓글 포함)
func (r *commentRepository) FindByPostIDWithReplies(postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment

	// 최상위 댓글 조회
	err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// 각 댓글에 대해 재귀적으로 대댓글 로드
	for i := range comments {
		r.loadReplies(&comments[i])
	}
	return comments, nil
}

func (r *commentRepository) Update(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}

func (r *commentRepository) HasReplies(commentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Comment{}).
		Where("parent_id = ?", commentID).
		Count(&count).Error
	return count > 0, err
}

func (r *commentRepository) loadReplies(comment *domain.Comment) {
	var replies []domain.Comment
	r.db.
		Where("parent_id = ?", comment.ID).
		Order("created_at ASC").
		Find(&replies)

	for i := range replies {
		r.loadReplies(&replies[i])
	}

	comment.Replies = replies
}
