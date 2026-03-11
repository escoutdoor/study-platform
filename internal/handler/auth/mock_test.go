package auth

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, in entity.User) (entity.Tokens, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Tokens), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, in entity.User) (entity.Tokens, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Tokens), args.Error(1)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (entity.Tokens, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(entity.Tokens), args.Error(1)
}
