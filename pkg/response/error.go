package response

import (
	"gorm-test/pkg/apperror"

	"github.com/gin-gonic/gin"
)

// Error는 에러를 저장하고 핸들러를 종료합니다
func Error(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// BadRequest는 400 에러를 반환합니다
func BadRequest(c *gin.Context, message string) {
	Error(c, apperror.BadRequest(message))
}

// NotFound는 404 에러를 반환합니다
func NotFound(c *gin.Context, resource string) {
	Error(c, apperror.NotFound(resource))
}

// Unauthorized는 401 에러를 반환합니다
func Unauthorized(c *gin.Context, message string) {
	Error(c, apperror.Unauthorized(message))
}
