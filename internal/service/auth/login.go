package auth

import (
	"context"
	"errors"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/apperror/code"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/hasher"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Login(ctx context.Context, in entity.User) (entity.Tokens, error) {
	user, err := s.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		appErr := new(apperror.Error)
		if errors.As(err, &appErr) && appErr.Code == code.UserNotFound {
			return entity.Tokens{}, apperror.ErrIncorrectCredentials
		}

		return entity.Tokens{}, errwrap.Wrap("get user by email from repository", err)
	}

	if match := hasher.CompareHashAndPassword(user.Password, in.Password); !match {
		return entity.Tokens{}, apperror.ErrIncorrectCredentials
	}

	roles, err := s.userRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("get user roles", err)
	}

	accessToken, err := s.tokenProvider.GenerateAccessToken(user.ID, roles)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate jwt access token", err)
	}
	refreshToken, err := s.tokenProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate jwt refresh token", err)
	}

	return entity.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
