package service

import (
	"context"
	"errors"
	"gorm-test/internal/config"
	"gorm-test/internal/domain"
	"gorm-test/internal/dto"
	"gorm-test/internal/repository"
	"gorm-test/middleware"
	"gorm-test/pkg/apperror"
	"gorm-test/pkg/metrics"
	"time"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("게시글을 찾을 수 없습니다")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type PostListResponse struct {
	Posts      []PostItem `json:"posts"`
	TotalCount int64      `json:"total_count"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

type PostItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	ViewCount int       `json:"view_count"`
	LikeCount int       `json:"like_count"`
	IsLiked   *bool     `json:"is_liked,omitempty"` // 로그인한 경우에만 포함
	IsMine    *bool     `json:"is_mine,omitempty"`  // 로그인한 경우에만 포함
}

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

func (s *PostService) Create(ctx context.Context, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	start := time.Now()

	// 중복 제목 검사
	//existing, _ := s.postRepo.FindByTitle(req.Title)
	//if existing != nil {
	//	return nil, apperror.Conflict("이미 같은 제목의 게시글이 있습니다")
	//}

	post := &domain.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: claims.UserID,
	}
	err := s.postRepo.Create(post)

	// DB 쿼리 시간 기록
	metrics.DBQueryDuration.WithLabelValues("insert", "posts").
		Observe(time.Since(start).Seconds())

	if err != nil {
		return nil, apperror.InternalError(err).WithDetail("게시글 생성 실패")
	}

	// 게시글 생성 카운터 증가
	metrics.PostsCreated.Inc()
	metrics.PostsTotal.Inc()
	return s.toResponse(post), nil
}

func (s *PostService) GetByID(id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundWithID("게시글", id)
		}
		return nil, apperror.InternalError(err).WithDetail("게시글 조회 중 오류")
	}

	_ = s.postRepo.IncrementViews(id)
	post.Views++

	return s.toResponse(post), nil
}

func (s *PostService) GetList(page, size int, search *dto.SearchParams, sort *dto.SortParams) ([]dto.PostListResponse, *dto.Meta, error) {
	pagination := dto.NewPagination(
		page,
		size,
		s.cfg.Pagination.DefaultSize,
		s.cfg.Pagination.MaxSize,
	)

	posts, total, err := s.postRepo.FindAll(pagination, search, sort)
	if err != nil {
		return nil, nil, err
	}

	list := make([]dto.PostListResponse, len(posts))
	for i, post := range posts {
		list[i] = dto.PostListResponse{
			ID:        post.ID,
			Title:     post.Title,
			Author:    post.Author.Username,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
		}
	}

	totalPages := int(total) / pagination.Size
	if int(total)%pagination.Size > 0 {
		totalPages++
	}

	meta := &dto.Meta{
		Page:       pagination.Page,
		Size:       pagination.Size,
		Total:      total,
		TotalPages: totalPages,
	}

	return list, meta, nil
}

// GetListByCursor 커서 기반 게시글 목록 조회
func (s *PostService) GetListByCursor(cursorStr string, size int) ([]dto.PostListResponse, *dto.CursorMeta, error) {
	if size < 1 {
		size = s.cfg.Pagination.DefaultSize
	}
	if size > s.cfg.Pagination.MaxSize {
		size = s.cfg.Pagination.MaxSize
	}

	// 커서 디코딩
	var cursor *dto.Cursor
	if cursorStr != "" {
		var err error
		cursor, err = dto.DecodeCursor(cursorStr)
		if err != nil {
			return nil, nil, errors.New("유효하지 않은 커서입니다")
		}
	}

	// 조회
	posts, err := s.postRepo.FindAllByCursor(cursor, size)
	if err != nil {
		return nil, nil, err
	}

	// 다음 페이지 존재 여부 확인
	hasMore := len(posts) > size
	if hasMore {
		posts = posts[:size] // 마지막 1개 제거
	}

	// DTO 변환
	list := make([]dto.PostListResponse, len(posts))
	for i, post := range posts {
		list[i] = dto.PostListResponse{
			ID:        post.ID,
			Title:     post.Title,
			Author:    post.Author.Username,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
		}
	}

	// 다음 커서 생성
	var nextCursor string
	if hasMore && len(posts) > 0 {
		last := posts[len(posts)-1]
		c := &dto.Cursor{ID: last.ID, CreatedAt: last.CreatedAt}
		nextCursor = c.Encode()
	}

	meta := &dto.CursorMeta{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	return list, meta, nil
}

func (s *PostService) Update(ctx context.Context, id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	// 권한 검사: 작성자 본인 또는 관리자만 수정 가능
	if post.AuthorID != claims.UserID && claims.Role != "admin" {
		return nil, ErrForbidden
	}
	// 작성자 확인
	//if post.AuthorID != userID {
	//	return nil, apperror.Forbidden("본인의 게시글만 수정할 수 있습니다")
	//}

	post.Title = req.Title
	post.Content = req.Content

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return s.toResponse(post), nil
}

func (s *PostService) Delete(ctx context.Context, id uint) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	// 권한 검사
	if post.AuthorID != claims.UserID && claims.Role != "admin" {
		return ErrForbidden
	}

	return s.postRepo.Delete(id)
}

func (s *PostService) toResponse(post *domain.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.Author.Username,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

//func (s *PostService) List(ctx context.Context, page, pageSize int) (*PostListResponse, error) {
//	posts, totalCount, err := s.postRepo.FindAll(ctx, page, pageSize)
//	if err != nil {
//		return nil, err
//	}
//
//	// 현재 사용자 (있을 수도 없을 수도)
//	claims, isLoggedIn := middleware.GetUserFromContext(ctx)
//
//	items := make([]PostItem, len(posts))
//	for i, post := range posts {
//		item := PostItem{
//			ID:        post.ID,
//			Title:     post.Title,
//			Author:    post.Author.Username,
//			CreatedAt: post.CreatedAt,
//			ViewCount: post.ViewCount,
//			LikeCount: post.LikeCount,
//		}
//
//		// 로그인한 경우에만 추가 정보 제공
//		if isLoggedIn {
//			isLiked := s.likeRepo.IsLiked(ctx, claims.UserID, post.ID)
//			isMine := post.AuthorID == claims.UserID
//			item.IsLiked = &isLiked
//			item.IsMine = &isMine
//		}
//
//		items[i] = item
//	}
//
//	return &PostListResponse{
//		Posts:      items,
//		TotalCount: totalCount,
//		Page:       page,
//		PageSize:   pageSize,
//	}, nil
//}
