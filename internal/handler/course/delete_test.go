package course

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

func TestDelete(t *testing.T) {
	type mockBehavior func(m *mockCourseService)

	tests := []struct {
		name string

		userID   int
		courseID string
		mockFn   mockBehavior
		noCtxID  bool

		wantCode int
		wantBody *httpresponse.ErrorResponse
	}{
		{
			name:     "success",
			userID:   1,
			courseID: "2",
			mockFn: func(m *mockCourseService) {
				m.On("Delete", mock.Anything, 2, 1).Return(nil).Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:     "invalid course id - not a number",
			userID:   1,
			courseID: "ivan",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - zero",
			userID:   1,
			courseID: "0",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - a negative number",
			userID:   1,
			courseID: "-1",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - empty",
			userID:   1,
			courseID: "0",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "missing user id in context",
			noCtxID:  true,
			courseID: "123",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: &httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:     "course not found",
			userID:   1,
			courseID: "123",
			mockFn: func(m *mockCourseService) {
				m.On("Delete", mock.Anything, 123, 1).
					Return(apperror.CourseNotFoundID(123)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: &httpresponse.ErrorResponse{
				Message: "course with id 123 was not found",
			},
		},
		{
			name:     "access denied - not course owner",
			userID:   2,
			courseID: "15",
			mockFn: func(m *mockCourseService) {
				m.On("Delete", mock.Anything, 15, 2).
					Return(apperror.CourseAccessDenied).Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: &httpresponse.ErrorResponse{
				Message: "only author can manage this course",
			},
		},
		{
			name:     "service internal error",
			userID:   1,
			courseID: "5",
			mockFn: func(m *mockCourseService) {
				m.On("Delete", mock.Anything, 5, 1).
					Return(errors.New("db connection lost")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockCourseService)
			cv := validator.New()
			h := &handler{service: svc, cv: cv}

			tt.mockFn(svc)

			r := httptest.NewRequest(http.MethodDelete, "/courses/"+tt.courseID, nil)
			r.SetPathValue("id", tt.courseID)

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.delete)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantBody != nil {
				wantJSON, err := json.Marshal(tt.wantBody)
				require.NoError(t, err)

				assert.JSONEq(t, string(wantJSON), w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}

			svc.AssertExpectations(t)
		})
	}
}
