package course

import (
	"context"
	"errors"
	"testing"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name string

		input  entity.Course
		mockFn func(m *mockCourseRepository)

		wantID  int
		wantErr error
	}{
		{
			name: "success",
			input: entity.Course{
				TeacherID:   1,
				Title:       "Good course",
				Description: "learn a lot there",
			},
			mockFn: func(m *mockCourseRepository) {
				in := entity.Course{
					TeacherID:   1,
					Title:       "Good course",
					Description: "learn a lot there",
				}
				m.On("Create", mock.Anything, in).Return(1, nil).Once()
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name: "repo error",
			input: entity.Course{
				TeacherID:   1,
				Title:       "Good course",
				Description: "learn a lot there",
			},
			mockFn: func(m *mockCourseRepository) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(0, errors.New("database error")).Once()
			},
			wantID:  0,
			wantErr: errors.New("create course in repo: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockCourseRepository)
			s := New(repo)

			tt.mockFn(repo)

			courseID, err := s.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantID, courseID)
			repo.AssertExpectations(t)
		})
	}
}
