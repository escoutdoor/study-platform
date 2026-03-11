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

func TestUnenroll(t *testing.T) {
	tests := []struct {
		name string

		userID   int
		courseID string
		noCtxID  bool
		mockFn   func(m *mockEnrollmentService)

		wantCode int
		wantBody *httpresponse.ErrorResponse
	}{
		{
			name:     "success",
			userID:   12,
			courseID: "13",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Unenroll", mock.Anything, 12, 13).
					Return(nil).Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:     "missing user id in ctx",
			noCtxID:  true,
			courseID: "3",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: &httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:     "invalid course id - not a number",
			userID:   1,
			courseID: "my-name",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusBadRequest,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - zero",
			userID:   1,
			courseID: "0",
			mockFn:   func(m *mockEnrollmentService) {},
			wantCode: http.StatusBadRequest,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "course not found",
			userID:   1,
			courseID: "404",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Unenroll", mock.Anything, 1, 404).
					Return(apperror.CourseNotFoundID(404)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course with id 404 was not found",
			},
		},
		{
			name:     "student not enrolled",
			userID:   1,
			courseID: "5",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Unenroll", mock.Anything, 1, 5).
					Return(apperror.StudentNotEnrolled).Once()
			},
			wantCode: http.StatusConflict,
			wantBody: &httpresponse.ErrorResponse{
				Message: "you are not enrolled in this course",
			},
		},
		{
			name:     "service returns internal error",
			userID:   1,
			courseID: "51",
			mockFn: func(m *mockEnrollmentService) {
				m.On("Unenroll", mock.Anything, 1, 51).
					Return(errors.New("something bad happened. we don't know")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &httpresponse.ErrorResponse{
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

			r := httptest.NewRequest(http.MethodDelete, "/students/me/courses/"+tt.courseID, nil)
			r.SetPathValue("courseId", tt.courseID)

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.unenroll).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantBody != nil {
				wantJSON, err := json.Marshal(tt.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantJSON), w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}

			enrollmentSvc.AssertExpectations(t)
		})
	}
}
