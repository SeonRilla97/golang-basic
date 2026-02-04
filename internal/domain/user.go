package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role은 사용자 역할을 나타냅니다.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User는 사용자 엔티티입니다.
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Email       string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"` // JSON 응답에서 제외
	Username    string         `gorm:"size:100;not null" json:"username"`
	Role        Role           `gorm:"size:20;default:user" json:"role"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName은 테이블 이름을 반환합니다.
func (User) TableName() string {
	return "users"
}

// BeforeCreate는 사용자 생성 전에 비밀번호를 해싱합니다. - 쓰지말것 :: 해싱된거 다시 해싱하거나, 유효성 검사 불가 -> 서비스 로직에 명시
func (u *User) BeforeCreateDeprecated(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(u.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}
