package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockCourseRepository struct {
	mock.Mock
}

func (m *mockCourseRepository) List(ctx context.Context) ([]entity.Course, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Course), args.Error(1)
}

func (m *mockCourseRepository) Get(ctx context.Context, courseID int) (entity.Course, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).(entity.Course), args.Error(1)
}

func (m *mockCourseRepository) Create(ctx context.Context, in entity.Course) (int, error) {
	args := m.Called(ctx, in)
	return args.Int(0), args.Error(1)
}

func (m *mockCourseRepository) Update(ctx context.Context, in entity.Course) (entity.Course, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Course), args.Error(1)
}

func (m *mockCourseRepository) Delete(ctx context.Context, courseID int) error {
	args := m.Called(ctx, courseID)
	return args.Error(0)
}
