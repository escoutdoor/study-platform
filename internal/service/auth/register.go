package auth

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/hasher"
	"github.com/escoutdoor/study-platform/internal/util/token"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Register(ctx context.Context, in entity.User) (entity.Tokens, error) {
	pw, err := hasher.HashPassword(in.Password)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("hash password", err)
	}
	in.Password = pw

	var userID int

	if txErr := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		userID, err = s.userRepo.Create(ctx, in)
		if err != nil {
			return errwrap.Wrap("create user in user repo", err)
		}

		if err = s.studentRepo.Create(ctx, entity.Student{UserID: userID}); err != nil {
			return errwrap.Wrap("create student in student repo", err)
		}

		return nil
	}); txErr != nil {
		return entity.Tokens{}, txErr
	}

	accessToken, err := s.tokenProvider.GenerateAccessToken(userID, []token.Role{token.RoleStudent})
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate jwt access token", err)
	}
	refreshToken, err := s.tokenProvider.GenerateRefreshToken(userID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate jwt refresh token", err)
	}

	return entity.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
