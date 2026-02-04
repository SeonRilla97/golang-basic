package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm-test/internal/config"
	"gorm-test/internal/domain"
	"gorm-test/internal/repository"
	"gorm-test/internal/router"
	"gorm-test/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostHandlerSuite 테스트 스위트
type PostHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

// SetupSuite 테스트 시작 전 1회 실행
func (s *PostHandlerSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// 테스트 DB 연결
	dsn := "host=localhost user=gouser password=gopassword dbname=godb_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)
	s.db = db

	// 마이그레이션
	db.AutoMigrate(&domain.Post{}, &domain.Comment{})

	// 의존성 주입
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultSize: 10,
			MaxSize:     100,
		},
	}

	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	postService := service.NewPostService(postRepo, cfg)
	commentService := service.NewCommentService(commentRepo, postRepo)
	postHandler := NewPostHandler(postService)
	commentHandler := NewCommentHandler(commentService)

	r := router.NewRouter(postHandler, commentHandler)
	s.router = r.Setup()
}

// TearDownSuite 테스트 종료 후 1회 실행
func (s *PostHandlerSuite) TearDownSuite() {
	// 테스트 테이블 삭제
	s.db.Migrator().DropTable(&domain.Comment{}, &domain.Post{})
}

// SetupTest 각 테스트 전 실행
func (s *PostHandlerSuite) SetupTest() {
	// 테이블 초기화
	s.db.Exec("TRUNCATE TABLE comments RESTART IDENTITY CASCADE")
	s.db.Exec("TRUNCATE TABLE posts RESTART IDENTITY CASCADE")
}

func TestPostHandlerSuite(t *testing.T) {
	suite.Run(t, new(PostHandlerSuite))
}

func (s *PostHandlerSuite) TestCreatePost() {
	// Given
	body := map[string]string{
		"title":   "테스트 게시글",
		"content": "테스트 내용입니다",
		"author":  "테스터",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	s.router.ServeHTTP(rec, req)

	// Then
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.True(s.T(), response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "테스트 게시글", data["title"])
	assert.Equal(s.T(), "테스터", data["author"])
}

func (s *PostHandlerSuite) TestGetPostList() {
	// Given: 게시글 3개 생성
	for i := 1; i <= 3; i++ {
		s.db.Create(&domain.Post{
			Title:   fmt.Sprintf("게시글 %d", i),
			Content: fmt.Sprintf("내용 %d", i),
			Author:  "테스터",
		})
	}

	req, _ := http.NewRequest("GET", "/api/v1/posts?page=1&size=10", nil)
	rec := httptest.NewRecorder()

	// When
	s.router.ServeHTTP(rec, req)

	// Then
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	data := response["data"].([]interface{})
	assert.Len(s.T(), data, 3)

	meta := response["meta"].(map[string]interface{})
	assert.Equal(s.T(), float64(3), meta["total"])
}

func (s *PostHandlerSuite) TestGetPostByID() {
	// Given
	post := &domain.Post{
		Title:   "조회 테스트",
		Content: "내용",
		Author:  "테스터",
	}
	s.db.Create(post)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
	rec := httptest.NewRecorder()

	// When
	s.router.ServeHTTP(rec, req)

	// Then
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "조회 테스트", data["title"])
	assert.Equal(s.T(), float64(1), data["views"]) // 조회수 증가
}

func (s *PostHandlerSuite) TestGetPostByID_NotFound() {
	req, _ := http.NewRequest("GET", "/api/v1/posts/999", nil)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *PostHandlerSuite) TestCreateComment() {
	// Given: 게시글 생성
	post := &domain.Post{Title: "테스트", Content: "내용", Author: "작성자"}
	s.db.Create(post)

	body := map[string]string{
		"content": "좋은 글이네요!",
		"author":  "댓글러",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST",
		fmt.Sprintf("/api/v1/posts/%d/comments", post.ID),
		bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	s.router.ServeHTTP(rec, req)

	// Then
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), "좋은 글이네요!", data["content"])
	assert.Equal(s.T(), float64(post.ID), data["post_id"])
}

func (s *PostHandlerSuite) TestCreateReply() {
	// Given: 게시글과 댓글 생성
	post := &domain.Post{Title: "테스트", Content: "내용", Author: "작성자"}
	s.db.Create(post)

	comment := &domain.Comment{PostID: post.ID, Content: "첫 댓글", Author: "댓글러"}
	s.db.Create(comment)

	body := map[string]interface{}{
		"content":   "대댓글입니다",
		"author":    "대댓글러",
		"parent_id": comment.ID,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST",
		fmt.Sprintf("/api/v1/posts/%d/comments", post.ID),
		bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	s.router.ServeHTTP(rec, req)

	// Then
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(s.T(), float64(comment.ID), data["parent_id"])
}
