package handler

import (
	"errors"
	"gorm-test/internal/dto"
	"gorm-test/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// Create 댓글 생성
// POST /api/v1/posts/:postId/comments
func (h *CommentHandler) Create(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 게시글 ID입니다"))
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	comment, err := h.commentService.Create(uint(postID), &req)
	if err != nil {
		if errors.Is(err, service.ErrPostNotExists) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "댓글 생성에 실패했습니다"))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(comment))
}

// GetByPostID 게시글의 댓글 목록 조회
// GET /api/v1/posts/:postId/comments
func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 게시글 ID입니다"))
		return
	}

	comments, err := h.commentService.GetByPostID(uint(postID))
	if err != nil {
		if errors.Is(err, service.ErrPostNotExists) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "댓글 조회에 실패했습니다"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(comments))
}

// Update 댓글 수정
// PUT /api/v1/posts/:postId/comments/:id
func (h *CommentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 댓글 ID입니다"))
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	comment, err := h.commentService.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "댓글 수정에 실패했습니다"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(comment))
}

// Delete 댓글 삭제
// DELETE /api/v1/posts/:postId/comments/:id
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 댓글 ID입니다"))
		return
	}

	err = h.commentService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "댓글 삭제에 실패했습니다"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
