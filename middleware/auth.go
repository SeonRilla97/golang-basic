package middleware

import (
	"errors"
	"gorm-test/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader는 토큰을 담는 헤더 이름입니다.
	AuthorizationHeader = "Authorization"
	// AuthorizationType은 인증 타입입니다.
	AuthorizationType = "Bearer"
	// ContextUserKey는 컨텍스트에 저장되는 사용자 정보 키입니다.
	ContextUserKey = "user"
)

var (
	ErrMissingToken  = errors.New("missing authorization token")
	ErrInvalidFormat = errors.New("invalid authorization format")
)

// AuthMiddleware는 JWT 인증 미들웨어입니다.
func AuthMiddleware(tokenService *auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Authorization 헤더 추출
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "인증 토큰이 필요합니다",
			})
			return
		}

		// 2. Bearer 토큰 형식 검증
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != AuthorizationType {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "잘못된 인증 형식입니다",
			})
			return
		}

		tokenString := parts[1]

		// 3. 토큰 검증
		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			handleTokenError(c, err)
			return
		}

		// 4. 컨텍스트에 사용자 정보 저장 [ 핸들러용 ]
		c.Set(ContextUserKey, claims)

		// Go Context에도 저장 [ 서비스용 ]
		ctx := SetUserToContext(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)

		// 5. 다음 핸들러로 진행
		c.Next()
	}
}

func handleTokenError(c *gin.Context, err error) {
	if errors.Is(err, auth.ErrExpiredToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "토큰이 만료되었습니다",
			"code":  "TOKEN_EXPIRED",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "유효하지 않은 토큰입니다",
	})
}

// GetCurrentUser는 컨텍스트에서 현재 사용자 정보를 추출합니다.
func GetCurrentUser(c *gin.Context) (*auth.CustomClaims, bool) {
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return nil, false
	}

	claims, ok := value.(*auth.CustomClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

// MustGetCurrentUser는 현재 사용자 정보를 추출합니다. 없으면 패닉입니다.
func MustGetCurrentUser(c *gin.Context) *auth.CustomClaims {
	claims, ok := GetCurrentUser(c)
	if !ok {
		panic("user not found in context - auth middleware not applied?")
	}
	return claims
}

// OptionalAuthMiddleware는 선택적 인증 미들웨어입니다.
// 토큰이 있으면 검증하고, 없으면 그냥 통과합니다.
func OptionalAuthMiddleware(tokenService *auth.TokenService, tokenStore auth.TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)

		// 토큰이 없으면 그냥 통과
		if authHeader == "" {
			c.Next()
			return
		}

		// 토큰이 있으면 검증 시도
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != AuthorizationType {
			// 형식이 잘못되어도 그냥 통과 (에러 반환 안 함)
			c.Next()
			return
		}

		tokenString := parts[1]

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			// 토큰이 유효하지 않아도 그냥 통과
			// 단, 만료된 토큰이면 클라이언트에게 알려줄 수 있음
			if errors.Is(err, auth.ErrExpiredToken) {
				c.Header("X-Token-Expired", "true")
			}
			c.Next()
			return
		}
		
		// 블랙리스트 확인
		if tokenStore != nil {
			tokenID := claims.RegisteredClaims.ID
			if tokenID == "" {
				tokenID = hashToken(tokenString)
			}

			blacklisted, err := tokenStore.IsBlacklisted(c.Request.Context(), tokenID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "토큰 검증 중 오류가 발생했습니다",
				})
				return
			}

			if blacklisted {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "토큰이 무효화되었습니다",
					"code":  "TOKEN_REVOKED",
				})
				return
			}
		}

		// 유효한 토큰이면 컨텍스트에 저장
		c.Set(ContextUserKey, claims)
		ctx := SetUserToContext(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
