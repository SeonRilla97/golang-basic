package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenService는 JWT 토큰 생성과 검증을 담당합니다.
type TokenService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	tokenStore    TokenStore
}

// NewTokenService는 TokenService를 생성합니다.
func NewTokenService(secretKey string, accessExpiry, refreshExpiry time.Duration, tokenStore TokenStore) *TokenService {
	return &TokenService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        "goboard-api",
		tokenStore:    tokenStore,
	}
}

// GenerateAccessToken은 액세스 토큰을 생성합니다.
func (s *TokenService) GenerateAccessToken(userID uint, email, username, role string) (string, error) {
	now := time.Now()
	tokenID := generateTokenID()

	claims := CustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken은 토큰을 검증하고 클레임을 반환합니다.
func (s *TokenService) ValidateToken(tokenString string) (*CustomClaims, error) {
	// Custom 헤더 사용 시 ParseWithClaims - 일반 파싱은 Parse()
	//수행 도중 CustomClaims.Validate 메소드를 실행함
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 알고리즘 검증
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secretKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}), // 허용하는 서명 알고리즘
		jwt.WithIssuer("goboard-api"),           // issuer 검증
		jwt.WithAudience("api.example.com"),     // audience 검증
		jwt.WithExpirationRequired(),            // exp 클레임 필수
		jwt.WithLeeway(5*time.Second),           // 시간 오차 허용 (서버 간 시간 차이 대응 - 방금 발급된 토큰이 "아직 유효하지 않음"으로 거부되는 것을 방지)
	)

	if err != nil {
		return nil, s.handleTokenError(err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// 토큰 에러 제어
func (s *TokenService) handleTokenError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		// 토큰 형식이 잘못됨 (점 3개가 아님 등)
		return fmt.Errorf("%w: malformed token", ErrInvalidToken)

	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// 서명이 유효하지 않음 (위조된 토큰)
		return fmt.Errorf("%w: invalid signature", ErrInvalidToken)

	case errors.Is(err, jwt.ErrTokenExpired):
		// 토큰이 만료됨
		return ErrExpiredToken

	case errors.Is(err, jwt.ErrTokenNotValidYet):
		// 아직 활성화되지 않은 토큰 (nbf 이전)
		return fmt.Errorf("%w: token not valid yet", ErrInvalidToken)

	default:
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
}

// GenerateRefreshToken은 리프레시 토큰을 생성합니다.
func (s *TokenService) GenerateRefreshToken(ctx context.Context, userID uint) (string, error) {
	now := time.Now()
	tokenID := generateTokenID()

	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ID:        generateTokenID(), // 고유 ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	// Redis에 저장
	if s.tokenStore != nil {
		if err := s.tokenStore.StoreRefreshToken(ctx, userID, tokenID, s.refreshExpiry); err != nil {
			return "", err
		}
	}

	return tokenString, nil
}

// generateTokenID는 토큰 고유 ID를 생성합니다.
func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetAccessExpiry는 액세스 토큰 만료 시간(초)을 반환합니다.
func (s *TokenService) GetAccessExpiry() int64 {
	return int64(s.accessExpiry.Seconds())
}

// ValidateRefreshToken은 리프레시 토큰을 검증합니다.
func (s *TokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return 0, s.handleTokenError(err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	// Subject에서 UserID 추출
	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, ErrInvalidToken
	}

	// Redis에서 유효성 확인
	if s.tokenStore != nil {
		valid, err := s.tokenStore.IsRefreshTokenValid(ctx, uint(userID), claims.ID)
		if err != nil {
			return 0, err
		}
		if !valid {
			// 이미 사용된 토큰 (Token Rotation으로 교체됨)
			return 0, ErrInvalidToken
		}
	}
	return uint(userID), nil
}

/*
RevokeAccessToken은 액세스 토큰을 무효화합니다.

	로그아웃: 사용자가 명시적으로 로그아웃
	비밀번호 변경: 기존 세션을 모두 끊어야 함
	계정 탈취 의심: 관리자가 강제로 세션 종료
	권한 변경: 권한이 낮아졌을 때 기존 토큰으로 접근 방지
*/

func (s *TokenService) RevokeAccessToken(ctx context.Context, tokenString string) error {
	if s.tokenStore == nil {
		return nil
	}

	// 토큰에서 ID와 만료 시간 추출
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.secretKey, nil
		},
	)

	if err != nil {
		// 이미 만료된 토큰이라도 파싱 시도
		// jwt.ErrTokenExpired인 경우에도 claims는 추출 가능
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return ErrInvalidToken
	}

	// 남은 만료 시간 계산
	expiry := time.Until(claims.ExpiresAt.Time)
	if expiry <= 0 {
		// 이미 만료됨 - 블랙리스트 추가 불필요
		return nil
	}

	// 토큰 ID로 블랙리스트 추가
	tokenID := claims.RegisteredClaims.ID
	if tokenID == "" {
		// 토큰 ID가 없으면 토큰 전체 해시 사용
		tokenID = hashToken(tokenString)
	}

	return s.tokenStore.AddToBlacklist(ctx, tokenID, expiry)
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:16]) // 앞 16바이트만 사용
}
