package student

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"github.com/escoutdoor/study-platform/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEnroll(t *testing.T) {
	tests := []struct {
		name string

		userID   int
		courseID string
		mockFn   func(m *mockEnrollmentService)
		noCtxID  bool

		wantCode int
		wantBody any
	}{
		{
			name:     "success",
			userID:   1,
			courseID: "15",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Enroll", mock.Anything, 1, 15).Return(nil).Once()
			},
			wantCode: http.StatusCreated,
			wantBody: map[string]string{"message": "successfully enrolled"},
		},
		{
			name:     "missing user id in context",
			noCtxID:  true,
			courseID: "3",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:     "invalid course id - not a number",
			userID:   1,
			courseID: "helloworld",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - zero",
			userID:   1,
			courseID: "0",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - negative",
			userID:   1,
			courseID: "-1111",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "course not found",
			userID:   1,
			courseID: "42",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Enroll", mock.Anything, 1, 42).
					Return(apperror.CourseNotFoundID(42)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "course with id 42 was not found",
			},
		},
		{
			name:     "already enrolled",
			userID:   1,
			courseID: "2",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Enroll", mock.Anything, 1, 2).
					Return(apperror.StudentAlreadyEnrolled).Once()
			},
			wantCode: http.StatusConflict,
			wantBody: httpresponse.ErrorResponse{
				Message: "you are already enrolled in this course",
			},
		},
		{
			name:     "internal error",
			userID:   1,
			courseID: "44",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Enroll", mock.Anything, 1, 44).
					Return(errors.New("something happened in repository layer")).Once()
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

			tt.mockFn(enrollmentSvc)

			r := httptest.NewRequest(http.MethodPost, "/students/me/courses/"+tt.courseID, nil)
			r.SetPathValue("courseId", tt.courseID)

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.enroll).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), w.Body.String())

			enrollmentSvc.AssertExpectations(t)
		})
	}
}
