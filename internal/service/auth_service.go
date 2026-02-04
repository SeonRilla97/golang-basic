package service

import (
	"context"
	"errors"
	"gorm-test/internal/auth"
	"gorm-test/internal/domain"
	"gorm-test/internal/dto"
	"gorm-test/internal/repository"
	"log/slog"
	"time"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}

type authService struct {
	userRepo        repository.UserRepository
	passwordService *auth.PasswordService
	tokenService    *auth.TokenService
}

func NewAuthService(
	userRepo repository.UserRepository,
	passwordService *auth.PasswordService,
	tokenService *auth.TokenService,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

func (s *authService) Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error) {
	// 1. 비밀번호 유효성 검사
	if err := s.passwordService.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// 2. 이메일 중복 확인
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// 3. 비밀번호 해싱
	hashedPassword, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// 4. 사용자 생성
	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
		Role:     domain.RoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	// 5. 응답 생성
	return &dto.SignupResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 사용자 조회
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		slog.Warn("login failed", "email", req.Email, "reasen", "user_not_found")
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. 비밀번호 검증
	if err := s.passwordService.Compare(user.Password, req.Password); err != nil {
		slog.Warn("login failed", "email", req.Email, "user_id", user.ID, "reason", "invalid_password")
		return nil, ErrInvalidCredentials
	}

	// 3. 토큰 생성
	accessToken, err := s.tokenService.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 4. 마지막 로그인 시간 업데이트
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user) // 에러는 무시 (로그인 성공에 영향 없음)
	slog.Info("login success", "user_id", user.ID, "email", user.Email)

	// 5. 응답 생성
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokenService.GetAccessExpiry(),
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     string(user.Role),
		},
	}, nil
}

/*
Token Rotation :: Access & Refresh Token 모두 재발급

Refresh Token 탈취 감지: 탈취된 토큰과 정상 토큰이 동시에 사용되면 감지 가능
만료 시간 연장: 활성 사용자는 계속 새 토큰을 받아 세션 유지
보안 강화: 오래된 Refresh Token이 자동으로 무효화
*/
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// 1. Refresh Token 검증
	userID, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 사용자 조회 (존재 여부 확인)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 3. 새 Access Token 생성
	newAccessToken, err := s.tokenService.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	// 4. 새 Refresh Token 생성 (Token Rotation)
	newRefreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokenService.GetAccessExpiry(),
	}, nil
}
