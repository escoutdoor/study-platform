package teacher

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

		mockFn func(m *mockTeacherRepository)

		wantTeachers []entity.Teacher
		wantErr      error
	}{
		{
			name: "success - multiple teachers",
			mockFn: func(m *mockTeacherRepository) {
				m.On("List", mock.Anything).Return([]entity.Teacher{
					{UserID: 1, FirstName: "Ivan", LastName: "Popov", Department: "CS", Email: "ivan@example.com"},
					{UserID: 2, FirstName: "Maria", LastName: "Piven", Department: "Math", Email: "maria@example.com"},
				}, nil).Once()
			},
			wantTeachers: []entity.Teacher{
				{UserID: 1, FirstName: "Ivan", LastName: "Popov", Department: "CS", Email: "ivan@example.com"},
				{UserID: 2, FirstName: "Maria", LastName: "Piven", Department: "Math", Email: "maria@example.com"},
			},
			wantErr: nil,
		},
		{
			name: "success - empty list",
			mockFn: func(m *mockTeacherRepository) {
				m.On("List", mock.Anything).Return([]entity.Teacher{}, nil).Once()
			},
			wantTeachers: []entity.Teacher{},
			wantErr:      nil,
		},
		{
			name: "repository error",
			mockFn: func(m *mockTeacherRepository) {
				m.On("List", mock.Anything).
					Return([]entity.Teacher(nil), errors.New("database error")).Once()
			},
			wantTeachers: nil,
			wantErr:      errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockTeacherRepository)
			s := New(repo)

			tt.mockFn(repo)

			teachers, err := s.List(context.Background())

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantTeachers, teachers)
			repo.AssertExpectations(t)
		})
	}
}
