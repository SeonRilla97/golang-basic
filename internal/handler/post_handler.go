package handler

import (
	"gorm-test/internal/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func (*PostHandler) Create(c *gin.Context) {
	// 요청 파싱

	// 유효성 검사

	// Service 호출

	// 응답 반환
}
