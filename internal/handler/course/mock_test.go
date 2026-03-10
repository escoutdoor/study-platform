package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockCourseService struct {
	mock.Mock
}

func (m *mockCourseService) Create(ctx context.Context, in entity.Course) (int, error) {
	args := m.Called(ctx, in)
	return args.Int(0), args.Error(1)
}

func (m *mockCourseService) Update(ctx context.Context, in entity.Course) (entity.Course, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Course), args.Error(1)
}

func (m *mockCourseService) Get(ctx context.Context, id int) (entity.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Course), args.Error(1)
}

func (m *mockCourseService) List(ctx context.Context) ([]entity.Course, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Course), args.Error(1)
}

func (m *mockCourseService) Delete(ctx context.Context, id, teacherID int) error {
	args := m.Called(ctx, id, teacherID)
	return args.Error(0)
}
