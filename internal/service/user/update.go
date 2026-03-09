package user

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/hasher"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Update(ctx context.Context, in entity.User) (entity.User, error) {
	if _, err := s.userRepo.GetByID(ctx, in.ID); err != nil {
		return entity.User{}, errwrap.Wrap("get user from repo by id", err)
	}

	pw, err := hasher.HashPassword(in.Password)
	if err != nil {
		return entity.User{}, errwrap.Wrap("hash password", err)
	}
	in.Password = pw

	updatedUser, err := s.userRepo.Update(ctx, in)
	if err != nil {
		return entity.User{}, errwrap.Wrap("update user in repo", err)
	}

	return updatedUser, nil
}
