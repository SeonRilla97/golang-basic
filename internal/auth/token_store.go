package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenStore는 토큰 저장소 인터페이스입니다.
type TokenStore interface {
	// StoreRefreshToken은 Refresh Token을 저장합니다.
	StoreRefreshToken(ctx context.Context, userID uint, tokenID string, expiry time.Duration) error

	// GetRefreshToken은 저장된 Refresh Token ID를 조회합니다.
	GetRefreshToken(ctx context.Context, userID uint) (string, error)

	// DeleteRefreshToken은 Refresh Token을 삭제합니다.
	DeleteRefreshToken(ctx context.Context, userID uint) error

	// IsRefreshTokenValid는 Refresh Token이 유효한지 확인합니다.
	IsRefreshTokenValid(ctx context.Context, userID uint, tokenID string) (bool, error)

	// AddToBlacklist는 토큰을 블랙리스트에 추가합니다.
	AddToBlacklist(ctx context.Context, tokenID string, expiry time.Duration) error

	// IsBlacklisted는 토큰이 블랙리스트에 있는지 확인합니다.
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

type RedisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(client *redis.Client) *RedisTokenStore {
	return &RedisTokenStore{client: client}
}

// 키 형식
func refreshTokenKey(userID uint) string {
	return fmt.Sprintf("refresh_token:%d", userID)
}

func blacklistKey(tokenID string) string {
	return fmt.Sprintf("blacklist:%s", tokenID)
}

func (s *RedisTokenStore) StoreRefreshToken(ctx context.Context, userID uint, tokenID string, expiry time.Duration) error {
	key := refreshTokenKey(userID)
	return s.client.Set(ctx, key, tokenID, expiry).Err()
}

func (s *RedisTokenStore) GetRefreshToken(ctx context.Context, userID uint) (string, error) {
	key := refreshTokenKey(userID)
	tokenID, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return tokenID, err
}

func (s *RedisTokenStore) DeleteRefreshToken(ctx context.Context, userID uint) error {
	key := refreshTokenKey(userID)
	return s.client.Del(ctx, key).Err()
}

func (s *RedisTokenStore) IsRefreshTokenValid(ctx context.Context, userID uint, tokenID string) (bool, error) {
	storedID, err := s.GetRefreshToken(ctx, userID)
	if err != nil {
		return false, err
	}
	return storedID == tokenID, nil
}

func (s *RedisTokenStore) AddToBlacklist(ctx context.Context, tokenID string, expiry time.Duration) error {
	key := blacklistKey(tokenID)
	return s.client.Set(ctx, key, "1", expiry).Err()
}

func (s *RedisTokenStore) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := blacklistKey(tokenID)
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
