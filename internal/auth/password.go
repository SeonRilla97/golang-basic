package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort        = errors.New("password must be at least 8 characters")
	ErrPasswordNoUpper         = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower         = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit         = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial       = errors.New("password must contain at least one special character")
	ErrPasswordMismatch  error = errors.New("password mismatch")
)

const (
	MinPasswordLength = 8
	DefaultCost       = 12 // 2^12 = 4096 iterations  || 10	~100ms	개발 환경 || 12	~400ms	일반 서비스 || 14	~1.5s	높은 보안

)

// PasswordService는 비밀번호 해싱과 검증을 담당합니다.
type PasswordService struct {
	cost int
}

// NewPasswordService는 PasswordService를 생성합니다.
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: DefaultCost,
	}
}

// Hash는 비밀번호를 해싱합니다.
func (s *PasswordService) Hash(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", ErrPasswordTooShort
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// Compare는 평문 비밀번호와 해시된 비밀번호를 비교합니다.
func (s *PasswordService) Compare(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return err
	}
	return nil
}

// ValidatePassword는 비밀번호 정책을 검사합니다.
func (s *PasswordService) ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}
