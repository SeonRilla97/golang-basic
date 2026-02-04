package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims는 JWT 토큰의 페이로드를 정의합니다.
type CustomClaims struct {
	UserID               uint   `json:"user_id"`
	Email                string `json:"email"`
	Username             string `json:"username"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // exp iat sub 등을 자동 상속
}

// Validate는 커스텀 검증 로직을 수행합니다.
// jwt.ClaimsValidator 인터페이스 구현
func (c CustomClaims) Validate() error {
	// Role 검증
	validRoles := map[string]bool{"user": true, "admin": true}
	if !validRoles[c.Role] {
		return fmt.Errorf("invalid role: %s", c.Role)
	}

	// UserID 검증
	if c.UserID == 0 {
		return errors.New("user_id is required")
	}

	// Email 형식 검증 (간단한 예시)
	if c.Email == "" {
		return errors.New("email is required")
	}

	return nil
}
