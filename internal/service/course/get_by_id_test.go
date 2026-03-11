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

func TestGet(t *testing.T) {
	tests := []struct {
		name string

		courseID int
		mockFn   func(m *mockCourseRepository)

		wantCourse entity.Course
		wantErr    error
	}{
		{
			name:     "success",
			courseID: 1,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:          1,
					TeacherID:   1,
					Title:       "Course #1",
					Description: "course #1",
				}, nil).Once()
			},
			wantCourse: entity.Course{
				ID:          1,
				TeacherID:   1,
				Title:       "Course #1",
				Description: "course #1",
			},
			wantErr: nil,
		},
		{
			name:     "course not found",
			courseID: 123,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 123).
					Return(entity.Course{}, apperror.CourseNotFoundID(123)).Once()
			},
			wantCourse: entity.Course{},
			wantErr:    apperror.CourseNotFoundID(123),
		},
		{
			name:     "repo returns error",
			courseID: 1,
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).
					Return(entity.Course{}, errors.New("db error")).Once()
			},
			wantCourse: entity.Course{},
			wantErr:    errors.New("get course from repo: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockCourseRepository)
			s := New(repo)

			tt.mockFn(repo)

			course, err := s.Get(context.Background(), tt.courseID)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCourse, course)
			repo.AssertExpectations(t)
		})
	}
}
