package teacher

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockTeacherRepository struct {
	mock.Mock
}

func (m *mockTeacherRepository) List(ctx context.Context) ([]entity.Teacher, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Teacher), args.Error(1)
}

func (m *mockTeacherRepository) GetByUserID(ctx context.Context, userID int) (entity.Teacher, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.Teacher), args.Error(1)
}

func (m *mockTeacherRepository) Create(ctx context.Context, in entity.Teacher) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *mockTeacherRepository) Update(ctx context.Context, in entity.Teacher) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}
