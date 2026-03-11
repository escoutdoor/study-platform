package student

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockStudentService struct {
	mock.Mock
}

func (m *mockStudentService) List(ctx context.Context) ([]entity.Student, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Student), args.Error(1)
}

func (m *mockStudentService) Get(ctx context.Context, userID int) (entity.Student, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.Student), args.Error(1)
}

func (m *mockStudentService) Update(ctx context.Context, in entity.Student) (entity.Student, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Student), args.Error(1)
}

type mockEnrollmentService struct {
	mock.Mock
}

func (m *mockEnrollmentService) Enroll(ctx context.Context, userID, courseID int) error {
	args := m.Called(ctx, userID, courseID)
	return args.Error(0)
}

func (m *mockEnrollmentService) Unenroll(ctx context.Context, userID, courseID int) error {
	args := m.Called(ctx, userID, courseID)
	return args.Error(0)
}
