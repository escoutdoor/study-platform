package course

import (
	"context"
	"errors"
	"testing"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name string

		courseID  int
		teacherID int
		mockFn    func(m *mockCourseRepository)

		wantErr error
	}{
		{
			name:      "success",
			courseID:  1,
			teacherID: 1,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:        1,
					TeacherID: 1,
				}, nil).Once()

				m.On("Delete", mock.Anything, 1).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:      "course not found",
			courseID:  123,
			teacherID: 1,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 123).
					Return(entity.Course{}, apperror.CourseNotFoundID(123)).Once()
			},
			wantErr: apperror.CourseNotFoundID(123),
		},
		{
			name:      "access denied - not owner",
			courseID:  1,
			teacherID: 2,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:        1,
					TeacherID: 1,
				}, nil).Once()
			},
			wantErr: apperror.CourseAccessDenied,
		},
		{
			name:      "get - repo error",
			courseID:  1,
			teacherID: 1,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).
					Return(entity.Course{}, errors.New("db error")).Once()
			},
			wantErr: errors.New("get course from repo: db error"),
		},
		{
			name:      "delete - repo error",
			courseID:  1,
			teacherID: 1,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:        1,
					TeacherID: 1,
				}, nil).Once()

				m.On("Delete", mock.Anything, 1).
					Return(errors.New("db error")).Once()
			},
			wantErr: errors.New("delete course in repo: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockCourseRepository)
			s := New(repo)

			tt.mockFn(repo)

			err := s.Delete(context.Background(), tt.courseID, tt.teacherID)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
