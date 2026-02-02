package service

import (
	"gorm-test/internal/domain"
	"gorm-test/internal/dto"
	"gorm-test/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func (s *PostService) Create(req *dto.CreatePostRequest) (*domain.Post, error) {
	// 1. 비즈니스 규칙 검증
	// 2. Repository 호출
	// 3. 결과 반환
}
