package auth

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (entity.Tokens, error) {
	claims, err := s.tokenProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("validate refresh token", err)
	}

	if _, err := s.userRepo.GetByID(ctx, claims.UserID); err != nil {
		return entity.Tokens{}, errwrap.Wrap("get user by id from repository", err)
	}

	roles, err := s.userRepo.GetRoles(ctx, claims.UserID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("get fresh user roles", err)
	}

	newAccessToken, err := s.tokenProvider.GenerateAccessToken(claims.UserID, roles)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate new access token", err)
	}
	newRefreshToken, err := s.tokenProvider.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate new refresh token", err)
	}

	return entity.Tokens{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}
