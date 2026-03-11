package teacher

import (
	"context"
	"errors"
	"testing"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name string

		input  entity.Teacher
		mockFn func(m *mockTeacherRepository)

		wantErr error
	}{
		{
			name: "success",
			input: entity.Teacher{
				UserID:     1,
				FirstName:  "Ivan",
				LastName:   "Popov",
				Department: "IT",
				Email:      "vanya@example.com",
			},
			mockFn: func(m *mockTeacherRepository) {
				m.On("Create", mock.Anything, entity.Teacher{
					UserID:     1,
					FirstName:  "Ivan",
					LastName:   "Popov",
					Department: "IT",
					Email:      "vanya@example.com",
				},
				).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "teacher already exists",
			input: entity.Teacher{
				UserID:     1,
				FirstName:  "Ivan",
				LastName:   "Popov",
				Department: "IT",
				Email:      "vanya@example.com",
			},
			mockFn: func(m *mockTeacherRepository) {
				m.On("Create", mock.Anything, entity.Teacher{
					UserID:     1,
					FirstName:  "Ivan",
					LastName:   "Popov",
					Department: "IT",
					Email:      "vanya@example.com",
				},
				).Return(apperror.TeacherAlreadyExists).Once()
			},
			wantErr: apperror.TeacherAlreadyExists,
		},
		{
			name: "teacher already exists",
			input: entity.Teacher{
				UserID:     1,
				FirstName:  "Ivan",
				LastName:   "Popov",
				Department: "IT",
				Email:      "vanya@example.com",
			},
			mockFn: func(m *mockTeacherRepository) {
				m.On("Create", mock.Anything, entity.Teacher{
					UserID:     1,
					FirstName:  "Ivan",
					LastName:   "Popov",
					Department: "IT",
					Email:      "vanya@example.com",
				},
				).Return(errors.New("database error")).Once()
			},
			wantErr: errwrap.Wrap("create teacher in repo", errors.New("database error")),
		},
	}

	for _, tt := range tests {
		repo := new(mockTeacherRepository)
		svc := New(repo)
		tt.mockFn(repo)

		err := svc.Create(context.Background(), tt.input)

		if tt.wantErr != nil {
			assert.ErrorContains(t, err, tt.wantErr.Error())
		} else {
			assert.NoError(t, err)
		}

		repo.AssertExpectations(t)
	}
}
