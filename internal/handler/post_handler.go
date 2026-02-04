package handler

import (
	"errors"
	"gorm-test/internal/dto"
	"gorm-test/internal/repository"
	"gorm-test/internal/service"
	"gorm-test/pkg/apperror"
	"gorm-test/pkg/response"

	"gorm-test/pkg/logger"
	"gorm-test/pkg/sentry"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.FromValidationErrors(err))
		return
	}

	post, err := h.postService.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(post))
}

func (h *PostHandler) GetByID(c *gin.Context) {
	log := logger.FromGin(c)
	id, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 ID 형식입니다")
		return
	}

	log.Debug("게시글 조회 시작", "post_id", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		response.Error(c, err)
		sentry.CaptureError(err)
		return
	}

	log.Info("게시글 조회 성공", "title", post.Title)
	response.Success(c, post)
}

func (h *PostHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 검색 파라미터 파싱
	search := &dto.SearchParams{
		Query:      c.Query("q"),
		SearchType: c.Query("type"),
	}

	sort := &dto.SortParams{
		Sort: c.Query("sort"),
	}

	posts, meta, err := h.postService.GetList(page, size, search, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "목록 조회에 실패했습니다"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(posts, meta))
}

// GetListByCursor 커서 기반 게시글 목록 조회
// GET /api/v1/posts/cursor?cursor=xxx&size=10
func (h *PostHandler) GetListByCursor(c *gin.Context) {
	cursor := c.Query("cursor")
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	posts, meta, err := h.postService.GetListByCursor(cursor, size)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_CURSOR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
		"meta":    meta,
	})
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	post, err := h.postService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "수정에 실패했습니다"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(post))
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	err = h.postService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "삭제에 실패했습니다"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PostHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "권한이 없습니다"})
	case errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}
