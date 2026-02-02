package handler

import (
	"errors"
	"gorm-test/internal/dto"
	"gorm-test/internal/service"
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	post, err := h.postService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "게시글 생성에 실패했습니다"))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(post))
}

func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "조회에 실패했습니다"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(post))
}

func (h *PostHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	posts, meta, err := h.postService.GetList(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "목록 조회에 실패했습니다"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(posts, meta))
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	post, err := h.postService.Update(uint(id), &req)
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	err = h.postService.Delete(uint(id))
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
