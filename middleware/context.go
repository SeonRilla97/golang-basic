package middleware

import (
	"context"
	"gorm-test/internal/auth"

	"github.com/gin-gonic/gin"
)

type contextKey string

const userContextKey contextKey = "user"

// SetUserToContext는 사용자 정보를 context.Context에 저장합니다.
func SetUserToContext(ctx context.Context, claims *auth.CustomClaims) context.Context {
	return context.WithValue(ctx, userContextKey, claims)
}

// GetUserFromContext는 context.Context에서 사용자 정보를 추출합니다.
func GetUserFromContext(ctx context.Context) (*auth.CustomClaims, bool) {
	claims, ok := ctx.Value(userContextKey).(*auth.CustomClaims)
	return claims, ok
}

// internal/middleware/context.go (추가)

// IsAuthenticated는 현재 요청이 인증되었는지 확인합니다.
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(ContextUserKey)
	return exists
}

// GetCurrentUserID는 현재 사용자 ID를 반환합니다.
// 인증되지 않은 경우 0을 반환합니다.
func GetCurrentUserID(c *gin.Context) uint {
	claims, ok := GetCurrentUser(c)
	if !ok {
		return 0
	}
	return claims.UserID
}

// GetCurrentUserOrNil은 현재 사용자를 반환합니다.
// 인증되지 않은 경우 nil을 반환합니다.
func GetCurrentUserOrNil(c *gin.Context) *auth.CustomClaims {
	claims, _ := GetCurrentUser(c)
	return claims
}
