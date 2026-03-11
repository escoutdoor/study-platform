package teacher

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockTeacherService struct {
	mock.Mock
}

func (m *mockTeacherService) List(ctx context.Context) ([]entity.Teacher, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Teacher), args.Error(1)
}

func (m *mockTeacherService) Get(ctx context.Context, userID int) (entity.Teacher, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.Teacher), args.Error(1)
}

func (m *mockTeacherService) Create(ctx context.Context, in entity.Teacher) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *mockTeacherService) Update(ctx context.Context, in entity.Teacher) (entity.Teacher, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Teacher), args.Error(1)
}
