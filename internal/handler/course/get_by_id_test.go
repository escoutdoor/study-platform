package course

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
	type mockBehavior func(m *mockCourseService)

	tests := []struct {
		name string

		courseID string
		mockFn   mockBehavior

		wantCode int
		wantBody any
	}{
		{
			name:     "success",
			courseID: "1",
			mockFn: func(m *mockCourseService) {
				m.On("Get", mock.Anything, 1).Return(entity.Course{
					ID:          1,
					TeacherID:   2,
					Title:       "Course #1",
					Description: "The best description",
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: getResponse{
				Course: courseResponse{
					ID:          1,
					TeacherID:   2,
					Title:       "Course #1",
					Description: "The best description",
				},
			},
		},
		{
			name:     "invalid course id - not a number",
			courseID: "hello",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - zero",
			courseID: "0",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - negative",
			courseID: "-1",
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "course not found",
			courseID: "122",
			mockFn: func(m *mockCourseService) {
				m.On("Get", mock.Anything, 122).
					Return(entity.Course{}, apperror.CourseNotFoundID(122)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "course with id 122 was not found",
			},
		},
		{
			name:     "service returns internal error",
			courseID: "5",
			mockFn: func(m *mockCourseService) {
				m.On("Get", mock.Anything, 5).
					Return(entity.Course{}, errors.New("db connection lost")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
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

			r := httptest.NewRequest(http.MethodGet, "/courses/"+tt.courseID, nil)
			r.SetPathValue("id", tt.courseID)

			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.get)
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
