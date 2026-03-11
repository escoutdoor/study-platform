package teacher

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

		userID int
		mockFn func(m *mockTeacherRepository)

		wantTeacher entity.Teacher
		wantErr     error
	}{
		{
			name:   "success",
			userID: 1,
			mockFn: func(m *mockTeacherRepository) {
				m.On("GetByUserID", mock.Anything, 1).Return(entity.Teacher{
					UserID:     1,
					FirstName:  "Ivan",
					LastName:   "Popov",
					Department: "IT",
					Email:      "ivan@example.com",
				}, nil).Once()
			},
			wantTeacher: entity.Teacher{
				UserID:     1,
				FirstName:  "Ivan",
				LastName:   "Popov",
				Department: "IT",
				Email:      "ivan@example.com",
			},
			wantErr: nil,
		},
		{
			name:   "teacher not found",
			userID: 404,
			mockFn: func(m *mockTeacherRepository) {
				m.On("GetByUserID", mock.Anything, 404).
					Return(entity.Teacher{}, apperror.TeacherNotFoundID(404)).Once()
			},
			wantTeacher: entity.Teacher{},
			wantErr:     apperror.TeacherNotFoundID(404),
		},
		{
			name:   "repo error",
			userID: 1,
			mockFn: func(m *mockTeacherRepository) {
				m.On("GetByUserID", mock.Anything, 1).
					Return(entity.Teacher{}, errors.New("database error")).Once()
			},
			wantTeacher: entity.Teacher{},
			wantErr:     errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockTeacherRepository)
			s := New(repo)

			tt.mockFn(repo)

			teacher, err := s.Get(context.Background(), tt.userID)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantTeacher, teacher)
			repo.AssertExpectations(t)
		})
	}
}
