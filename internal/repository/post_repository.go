package repository

import (
	"errors"
	"gorm-test/internal/domain"
	"gorm-test/internal/dto"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("board is not exist")
)

// PostRepository 게시글 저장소 인터페이스
type PostRepository interface {
	Create(post *domain.Post) error
	FindByID(id uint) (*domain.Post, error)
	FindAll(pagination *dto.Pagination, search *dto.SearchParams, sort *dto.SortParams) ([]domain.Post, int64, error)
	Update(post *domain.Post) error
	Delete(id uint) error
	IncrementViews(id uint) error
	FindAllByCursor(cursor *dto.Cursor, limit int) ([]domain.Post, error)
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

func (r *postRepository) FindAll(pagination *dto.Pagination, search *dto.SearchParams, sort *dto.SortParams) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	query := r.db.Model(&domain.Post{})

	// 검색 조건 적용
	/**
	LIKE 검색은 간단하지만 한계가 있습니다.

	인덱스를 타지 않아 대량 데이터에서 느림
	형태소 분석 없음 (한글 부분 검색 어려움)
	검색 결과 정렬이 단순함

	==> 대규모 서비스에서는 Elasticsearch 같은 전문 검색 엔진을 사용합니다. 하지만 우리 게시판 규모에서는 LIKE로 충분합니다.
	*/
	if search != nil && search.Query != "" {
		searchQuery := "%" + search.Query + "%"
		switch search.GetSearchType() {
		case dto.SearchTypeTitle:
			query = query.Where("title ILIKE ?", searchQuery) // 대소문자를 구분하지 않는 LIKE (Postgresql)
		case dto.SearchTypeContent:
			query = query.Where("content ILIKE ?", searchQuery)
		default: // all
			query = query.Where("title ILIKE ? OR content ILIKE ?", searchQuery, searchQuery)
		}
	}

	// 전체 개수 조회
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 정렬 조건 적용
	orderStr := "created_at DESC"
	if sort != nil {
		orderStr = sort.ToOrderString()
	}

	// 페이징 적용하여 조회
	err := r.db.
		Order(orderStr).
		Offset(pagination.Offset()).
		Limit(pagination.Size).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// 구현체에 추가
func (r *postRepository) FindAllByCursor(cursor *dto.Cursor, limit int) ([]domain.Post, error) {
	var posts []domain.Post

	query := r.db.Order("created_at DESC, id DESC")

	// 커서가 있으면 조건 추가
	if cursor != nil {
		query = query.Where(
			"(created_at < ?) OR (created_at = ? AND id < ?)",
			cursor.CreatedAt, cursor.CreatedAt, cursor.ID,
		)
	}

	err := query.Limit(limit + 1).Find(&posts).Error // 1개 더 조회해서 다음 페이지 존재 확인
	if err != nil {
		return nil, err
	}

	return posts, nil
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
