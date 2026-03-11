package student

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"github.com/escoutdoor/study-platform/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tests := []struct {
		name string

		mockFn func(m *mockStudentService)

		wantCode int
		wantBody any
	}{
		{
			name: "success - multiple students",
			mockFn: func(m *mockStudentService) {
				m.On("List", mock.Anything).Return([]entity.Student{
					{
						UserID:    1,
						FirstName: "Ivan",
						LastName:  "Popov",
						Email:     "ivan@example.com",
					},
					{
						UserID:    2,
						FirstName: "Anna",
						LastName:  "Koval",
						Email:     "anna@example.com",
					},
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Students: []studentResponse{
					{UserID: 1, FirstName: "Ivan", LastName: "Popov", Email: "ivan@example.com"},
					{UserID: 2, FirstName: "Anna", LastName: "Koval", Email: "anna@example.com"},
				},
			},
		},
		{
			name: "success - empty list",
			mockFn: func(m *mockStudentService) {
				m.On("List", mock.Anything).Return([]entity.Student{}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Students: []studentResponse{},
			},
		},
		{
			name: "service returns internal error",
			mockFn: func(m *mockStudentService) {
				m.On("List", mock.Anything).
					Return([]entity.Student(nil), errors.New("db connection lost")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			studentSvc := new(mockStudentService)
			enrollmentSvc := new(mockEnrollmentService)
			h := &handler{studentService: studentSvc, enrollmentService: enrollmentSvc, cv: validator.New()}

			tt.mockFn(studentSvc)

			r := httptest.NewRequest(http.MethodGet, "/students", nil)
			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.list).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), w.Body.String())

			studentSvc.AssertExpectations(t)
		})
	}
}
