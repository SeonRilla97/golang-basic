package service

import (
	"errors"
	"gorm-test/internal/domain"
	"gorm-test/internal/dto"
	"gorm-test/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrCommentNotFound = errors.New("댓글을 찾을 수 없습니다")
	ErrPostNotExists   = errors.New("게시글이 존재하지 않습니다")
)

const MaxReplyDepth = 3 // 최대 3단계까지

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// Create 댓글 생성
func (s *CommentService) Create(postID uint, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	// 게시글 존재 확인
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotExists
		}
		return nil, err
	}

	// 부모 댓글 확인 (대댓글인 경우)
	if req.ParentID != nil {
		parent, err := s.commentRepo.FindByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("부모 댓글을 찾을 수 없습니다")
		}
		// 같은 게시글인지 확인
		if parent.PostID != postID {
			return nil, errors.New("부모 댓글이 다른 게시글에 속합니다")
		}

		// 깊이 확인
		depth := s.getCommentDepth(parent)
		if depth >= MaxReplyDepth {
			return nil, errors.New("더 이상 대댓글을 작성할 수 없습니다")
		}

	}

	comment := &domain.Comment{
		PostID:   postID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Author:   req.Author,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return s.toResponse(comment), nil
}

// GetByPostID 게시글의 댓글 목록 조회
func (s *CommentService) GetByPostID(postID uint) ([]*dto.CommentResponse, error) {
	// 게시글 존재 확인
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotExists
		}
		return nil, err
	}

	comments, err := s.commentRepo.FindByPostIDWithReplies(postID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CommentResponse, len(comments))
	for i, comment := range comments {
		result[i] = s.toResponseWithReplies(&comment)
	}

	return result, nil
}

// Update 댓글 수정
func (s *CommentService) Update(id uint, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}

	comment.Content = req.Content

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return s.toResponse(comment), nil
}

// Delete 댓글 삭제 (내용만 삭제)
func (s *CommentService) Delete(id uint) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}

	// 대댓글이 있으면 내용만 변경
	hasReplies, _ := s.commentRepo.HasReplies(id)
	if hasReplies {
		comment.Content = "[삭제된 댓글입니다]"
		comment.Author = ""
		return s.commentRepo.Update(comment)
	}

	return s.commentRepo.Delete(id)
}
func (s *CommentService) toResponse(comment *domain.Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Author:    comment.Author,
		CreatedAt: comment.CreatedAt,
	}
}

// getCommentDepth 댓글 깊이 계산
func (s *CommentService) getCommentDepth(comment *domain.Comment) int {
	depth := 1
	current := comment

	for current.ParentID != nil {
		parent, err := s.commentRepo.FindByID(*current.ParentID)
		if err != nil {
			break
		}
		current = parent
		depth++
	}

	return depth
}

// toResponseWithReplies 대댓글 포함 변환
func (s *CommentService) toResponseWithReplies(comment *domain.Comment) *dto.CommentResponse {
	resp := &dto.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Author:    comment.Author,
		CreatedAt: comment.CreatedAt,
	}

	if len(comment.Replies) > 0 {
		resp.Replies = make([]*dto.CommentResponse, len(comment.Replies))
		for i, reply := range comment.Replies {
			resp.Replies[i] = s.toResponseWithReplies(&reply)
		}
	}

	return resp
}
