package response

import (
	"gorm-test/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success는 200 OK와 함께 데이터를 반환합니다
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, dto.SuccessResponse(data))
}

// Created는 201 Created와 함께 데이터를 반환합니다
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, dto.SuccessResponse(data))
}

// SuccessWithMeta는 200 OK와 함께 데이터 + 페이징 메타를 반환합니다
func SuccessWithMeta(c *gin.Context, data any, meta *dto.Meta) {
	c.JSON(http.StatusOK, dto.SuccessWithMeta(data, meta))
}

// NoContent는 204 No Content를 반환합니다
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
