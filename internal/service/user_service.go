package service

import "context"

func (s *userService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 현재 비밀번호 확인
	if err := s.passwordService.Compare(user.Password, currentPassword); err != nil {
		return ErrInvalidCredentials
	}

	// 새 비밀번호 해싱
	hashedPassword, err := s.passwordService.Hash(newPassword)
	if err != nil {
		return err
	}

	// 비밀번호 업데이트
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// 모든 토큰 무효화
	_ = s.tokenService.RevokeRefreshToken(ctx, userID)
	// Access Token은 짧은 만료 시간이므로 블랙리스트에 추가하지 않아도 됨

	return nil
}
