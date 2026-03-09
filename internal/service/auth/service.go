package auth

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/token"
	"github.com/escoutdoor/study-platform/pkg/database"
)

type Service struct {
	userRepo    userRepository
	studentRepo studentRepository
	txManager   database.TxManager

	tokenProvider tokenProvider
}

func New(
	userRepo userRepository,
	studentRepo studentRepository,
	txManager database.TxManager,
	tokenProvider tokenProvider,
) *Service {
	return &Service{
		userRepo:      userRepo,
		studentRepo:   studentRepo,
		txManager:     txManager,
		tokenProvider: tokenProvider,
	}
}

type studentRepository interface {
	Create(ctx context.Context, in entity.Student) error
}

type userRepository interface {
	Create(ctx context.Context, in entity.User) (int, error)

	GetByID(ctx context.Context, userID int) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetRoles(ctx context.Context, userID int) ([]token.Role, error)
}

type tokenProvider interface {
	ValidateAccessToken(accessToken string) (token.AccessTokenClaims, error)
	ValidateRefreshToken(refreshToken string) (token.RefreshTokenClaims, error)

	GenerateAccessToken(userID int, roles []token.Role) (string, error)
	GenerateRefreshToken(userID int) (string, error)
}
