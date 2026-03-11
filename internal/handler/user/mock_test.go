package user

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Update(ctx context.Context, in entity.User) (entity.User, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *mockUserService) Delete(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
