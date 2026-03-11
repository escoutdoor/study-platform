package student

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"github.com/escoutdoor/study-platform/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name string

		userID  int
		body    any
		mockFn  func(m *mockStudentService)
		noCtxID bool

		wantCode int
		wantBody any
	}{
		{
			name:   "success",
			userID: 1,
			body:   updateRequest{},
			mockFn: func(m *mockStudentService) {
				in := entity.Student{UserID: 1}
				m.On("Update", mock.Anything, in).Return(entity.Student{
					UserID:    1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: updateResponse{
				Student: studentResponse{
					UserID:    1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
				},
			},
		},
		{
			name:     "invalid json",
			userID:   1,
			body:     "not a json{{{",
			mockFn:   func(m *mockStudentService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name:     "missing user id in context",
			noCtxID:  true,
			body:     updateRequest{},
			mockFn:   func(m *mockStudentService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "student not found",
			userID: 999,
			body:   updateRequest{},
			mockFn: func(m *mockStudentService) {
				in := entity.Student{UserID: 999}
				m.On("Update", mock.Anything, in).
					Return(entity.Student{}, apperror.StudentNotFoundID(999)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "student with id 999 was not found",
			},
		},
		{
			name:   "service returns internal error",
			userID: 1,
			body:   updateRequest{},
			mockFn: func(m *mockStudentService) {
				in := entity.Student{UserID: 1}
				m.On("Update", mock.Anything, in).
					Return(entity.Student{}, errors.New("db connection lost")).Once()
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

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodPut, "/students/me", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.update).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), w.Body.String())

			studentSvc.AssertExpectations(t)
		})
	}
}
