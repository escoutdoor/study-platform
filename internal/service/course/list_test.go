package course

import (
	"context"
	"errors"
	"testing"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestList(t *testing.T) {
	tests := []struct {
		name string

		mockFn func(m *mockCourseRepository)

		wantCourses []entity.Course
		wantErr     error
	}{
		{
			name: "success - multiple courses",
			mockFn: func(m *mockCourseRepository) {
				m.On("List", mock.Anything).Return([]entity.Course{
					{ID: 1, TeacherID: 1, Title: "go basics", Description: "course about go"},
					{ID: 2, TeacherID: 2, Title: "Advanced SQL", Description: "deep dive into sql"},
				}, nil).Once()
			},
			wantCourses: []entity.Course{
				{ID: 1, TeacherID: 1, Title: "go basics", Description: "course about go"},
				{ID: 2, TeacherID: 2, Title: "Advanced SQL", Description: "deep dive into sql"},
			},
			wantErr: nil,
		},
		{
			name: "success - empty list",
			mockFn: func(m *mockCourseRepository) {
				m.On("List", mock.Anything).Return([]entity.Course{}, nil).Once()
			},
			wantCourses: []entity.Course{},
			wantErr:     nil,
		},
		{
			name: "repo returns error",
			mockFn: func(m *mockCourseRepository) {
				m.On("List", mock.Anything).
					Return([]entity.Course(nil), errors.New("db error")).Once()
			},
			wantCourses: nil,
			wantErr:     errors.New("get list of courses from repo: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockCourseRepository)
			s := New(repo)

			tt.mockFn(repo)

			courses, err := s.List(context.Background())

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCourses, courses)
			repo.AssertExpectations(t)
		})
	}
}
