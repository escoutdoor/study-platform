package course

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
	type mockBehavior func(m *mockCourseService)

	tests := []struct {
		name string

		mockFn mockBehavior

		wantCode int
		wantBody any
	}{
		{
			name: "success - multiple courses",
			mockFn: func(m *mockCourseService) {
				m.On("List", mock.Anything).Return([]entity.Course{
					{
						ID:          1,
						TeacherID:   2,
						Title:       "The first course",
						Description: "the best description ever",
					},
					{
						ID:          2,
						TeacherID:   4,
						Title:       "How to center div? -_-",
						Description: "the best way to do it",
					},
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Courses: []courseResponse{
					{
						ID:          1,
						TeacherID:   2,
						Title:       "The first course",
						Description: "the best description ever",
					},
					{
						ID:          2,
						TeacherID:   4,
						Title:       "How to center div? -_-",
						Description: "the best way to do it",
					},
				},
			},
		},
		{
			name: "success - 1 course",
			mockFn: func(m *mockCourseService) {
				m.On("List", mock.Anything).Return([]entity.Course{
					{
						ID:          1,
						TeacherID:   3,
						Title:       "How to create golang application",
						Description: "This course will teach you a lot",
					},
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Courses: []courseResponse{
					{
						ID:          1,
						TeacherID:   3,
						Title:       "How to create golang application",
						Description: "This course will teach you a lot",
					},
				},
			},
		},
		{
			name: "success - empty list",
			mockFn: func(m *mockCourseService) {
				m.On("List", mock.Anything).Return([]entity.Course{}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Courses: []courseResponse{},
			},
		},
		{
			name: "service returns internal error",
			mockFn: func(m *mockCourseService) {
				m.On("List", mock.Anything).
					Return([]entity.Course(nil), errors.New("db connection lost")).Once()
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

			r := httptest.NewRequest(http.MethodGet, "/courses", nil)
			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.list)
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
