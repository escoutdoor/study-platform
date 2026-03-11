package student

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"github.com/escoutdoor/study-platform/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string

		studentID string
		mockFn    func(m *mockStudentService)

		wantCode int
		wantBody any
	}{
		{
			name:      "success",
			studentID: "1",
			mockFn: func(m *mockStudentService) {
				m.On("Get", mock.Anything, 1).Return(
					entity.Student{
						UserID:    1,
						FirstName: "Hello",
						LastName:  "Super Hello",
						Email:     "example@example.com",
					},
					nil,
				).Once()
			},
			wantCode: http.StatusOK,
			wantBody: getResponse{
				Student: studentResponse{
					UserID:    1,
					FirstName: "Hello",
					LastName:  "Super Hello",
					Email:     "example@example.com",
				},
			},
		},
		{
			name:      "invalid student id - a negative number",
			studentID: "-363",
			mockFn:    func(m *mockStudentService) {},
			wantCode:  http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "student id must be a positive integer",
			},
		},
		{
			name:      "invalid student id - not a number",
			studentID: "popov",
			mockFn:    func(m *mockStudentService) {},
			wantCode:  http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "student id must be a positive integer",
			},
		},
		{
			name:      "invalid student id - zero",
			studentID: "0",
			mockFn:    func(m *mockStudentService) {},
			wantCode:  http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "student id must be a positive integer",
			},
		},
		{
			name:      "service returns internal error",
			studentID: "1",
			mockFn: func(m *mockStudentService) {
				m.On("Get", mock.Anything, 1).
					Return(entity.Student{}, errors.New("db connection lost")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:      "not found",
			studentID: "98333",
			mockFn: func(m *mockStudentService) {
				m.On("Get", mock.Anything, 98333).Return(entity.Student{}, apperror.StudentNotFoundID(98333)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "student with id 98333 was not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			studentService := new(mockStudentService)
			enrollmentService := new(mockEnrollmentService)
			h := &handler{
				studentService:    studentService,
				enrollmentService: enrollmentService,
				cv:                validator.New(),
			}

			tt.mockFn(studentService)

			r := httptest.NewRequest(http.MethodGet, "/students/"+tt.studentID, nil)
			r.SetPathValue("id", tt.studentID)

			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.get)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantBodyJson, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantBodyJson), w.Body.String())

			studentService.AssertExpectations(t)
		})
	}
}
