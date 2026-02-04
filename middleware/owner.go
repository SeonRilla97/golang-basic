package middleware

import (
	"gorm-test/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RequirePostOwner는 게시글 소유자만 허용합니다.
//
//	복잡한 권한 검사는 서비스 레이어에서 하는 것이 더 깔끔할 수 있습니다.
func RequirePostOwner(postRepo repository.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetCurrentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "인증이 필요합니다",
			})
			return
		}

		// 관리자는 통과
		if claims.Role == "admin" {
			c.Next()
			return
		}

		// 게시글 ID 추출
		postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "잘못된 게시글 ID",
			})
			return
		}

		// 게시글 조회
		post, err := postRepo.FindByID(c.Request.Context(), uint(postID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "게시글을 찾을 수 없습니다",
			})
			return
		}

		// 소유자 검증
		if post.AuthorID != claims.UserID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "게시글에 대한 권한이 없습니다",
			})
			return
		}

		c.Next()
	}
}
