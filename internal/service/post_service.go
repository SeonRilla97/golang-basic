package service

import (
	"errors"
	"gorm-test/internal/config"
	"gorm-test/internal/domain"
	"gorm-test/internal/dto"
	"gorm-test/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("게시글을 찾을 수 없습니다")
)

type PostService struct {
	postRepo repository.PostRepository
	cfg      *config.Config
}

func NewPostService(postRepo repository.PostRepository, cfg *config.Config) *PostService {
	return &PostService{
		postRepo: postRepo,
		cfg:      cfg,
	}
}

func (s *PostService) Create(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	post := &domain.Post{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	return s.toResponse(post), nil
}

func (s *PostService) GetByID(id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	_ = s.postRepo.IncrementViews(id)
	post.Views++

	return s.toResponse(post), nil
}

func (s *PostService) GetList(page, size int) ([]dto.PostListResponse, *dto.Meta, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = s.cfg.Pagination.DefaultSize
	}
	if size > s.cfg.Pagination.MaxSize {
		size = s.cfg.Pagination.MaxSize
	}

	offset := (page - 1) * size

	posts, total, err := s.postRepo.FindAll(offset, size)
	if err != nil {
		return nil, nil, err
	}

	list := make([]dto.PostListResponse, len(posts))
	for i, post := range posts {
		list[i] = dto.PostListResponse{
			ID:        post.ID,
			Title:     post.Title,
			Author:    post.Author,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
		}
	}

	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	meta := &dto.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}

	return list, meta, nil
}

func (s *PostService) Update(id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return s.toResponse(post), nil
}

func (s *PostService) Delete(id uint) error {
	_, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	return s.postRepo.Delete(id)
}

func (s *PostService) toResponse(post *domain.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.Author,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
