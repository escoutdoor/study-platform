package user

import (
	"context"

	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Delete(ctx context.Context, userID int) error {
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return errwrap.Wrap("get user from repo by id", err)
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errwrap.Wrap("delete user in repo", err)
	}

	return nil
}
