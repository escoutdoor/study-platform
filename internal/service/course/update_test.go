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

func TestUpdate(t *testing.T) {
	tests := []struct {
		name string

		input  entity.Course
		mockFn func(m *mockCourseRepository)

		wantCourse entity.Course
		wantErr    error
	}{
		{
			name: "success",
			input: entity.Course{
				ID:          1,
				TeacherID:   1,
				Title:       "New Course Availabel!",
				Description: "Updated description",
			},
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:        1,
					TeacherID: 1,
				}, nil).Once()

				m.On("Update", mock.Anything, entity.Course{
					ID:          1,
					TeacherID:   1,
					Title:       "New Course Availabel!",
					Description: "Updated description",
				}).Return(entity.Course{
					ID:          1,
					TeacherID:   1,
					Title:       "New Course Availabel!",
					Description: "Updated description",
				}, nil).Once()
			},
			wantCourse: entity.Course{
				ID:          1,
				TeacherID:   1,
				Title:       "New Course Availabel!",
				Description: "Updated description",
			},
			wantErr: nil,
		},
		{
			name: "course not found",
			input: entity.Course{
				ID:        404,
				TeacherID: 1,
			},
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 404).
					Return(entity.Course{}, apperror.CourseNotFoundID(404)).Once()
			},
			wantCourse: entity.Course{},
			wantErr:    apperror.CourseNotFoundID(404),
		},
		{
			name: "access denied - not owner",
			input: entity.Course{
				ID:        1,
				TeacherID: 2,
			},
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:        1,
					TeacherID: 1,
				}, nil).Once()
			},
			wantCourse: entity.Course{},
			wantErr:    apperror.CourseAccessDenied,
		},
		{
			name: "get - repo error",
			input: entity.Course{
				ID:        1,
				TeacherID: 1,
			},
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).
					Return(entity.Course{}, errors.New("db error")).Once()
			},
			wantCourse: entity.Course{},
			wantErr:    errors.New("get course from repo by id: db error"),
		},
		{
			name: "update - repo error",
			input: entity.Course{
				ID:          1,
				TeacherID:   1,
				Title:       "Updated",
				Description: "Updated",
			},
			mockFn: func(m *mockCourseRepository) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:        1,
					TeacherID: 1,
				}, nil).Once()

				m.On("Update", mock.Anything, mock.Anything).
					Return(entity.Course{}, errors.New("db error")).Once()
			},
			wantCourse: entity.Course{},
			wantErr:    errors.New("update course in repo: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockCourseRepository)
			s := New(repo)

			tt.mockFn(repo)

			course, err := s.Update(context.Background(), tt.input)

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
