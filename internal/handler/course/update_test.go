package course

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
	type mockBehavior func(m *mockCourseService)

	tests := []struct {
		name string

		userID      int
		courseID    string
		noCtxID     bool
		mockFn      mockBehavior
		requestBody any

		wantCode int
		wantBody any
	}{
		{
			name:     "success",
			userID:   1,
			courseID: "2",
			mockFn: func(m *mockCourseService) {
				m.On("Update", mock.Anything, entity.Course{
					ID:          2,
					TeacherID:   1,
					Title:       "Starter course",
					Description: "starter course description",
				}).Return(entity.Course{
					ID:          2,
					TeacherID:   1,
					Title:       "Starter course",
					Description: "starter course description",
				}, nil).Once()
			},
			requestBody: updateRequest{
				Title:       "Starter course",
				Description: "starter course description",
			},
			wantCode: http.StatusOK,
			wantBody: updateResponse{
				Course: courseResponse{
					ID:          2,
					TeacherID:   1,
					Title:       "Starter course",
					Description: "starter course description",
				},
			},
		},
		{
			name:        "invalid json",
			userID:      1,
			courseID:    "4",
			requestBody: "not a json{{{",
			mockFn:      func(m *mockCourseService) {},
			wantCode:    http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name:     "validation error - empty title",
			userID:   1,
			courseID: "3",
			requestBody: updateRequest{
				Title:       "",
				Description: "some cool description",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "title",
						Message: "This field is required",
					},
				},
			},
		},
		{
			name:     "validation error - title too short",
			userID:   1,
			courseID: "51",
			requestBody: updateRequest{
				Title:       "XD",
				Description: "some funny description",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "title",
						Message: "This field must be at least 3 characters long",
					},
				},
			},
		},
		{
			name:     "validation error - title too long",
			userID:   1,
			courseID: "4",
			requestBody: updateRequest{
				Title:       "This title is way too long. something something something something something something something",
				Description: "some fascinated description",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "title",
						Message: "This field must be at most 50 characters long",
					},
				},
			},
		},
		{
			name:     "validation error - empty description",
			userID:   1,
			courseID: "5",
			requestBody: updateRequest{
				Title:       "Course you would like to enroll in",
				Description: "",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "description",
						Message: "This field is required",
					},
				},
			},
		},
		{
			name:     "validation error - description too short",
			userID:   1,
			courseID: "44",
			requestBody: updateRequest{
				Title:       "Go information",
				Description: "ab",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "description",
						Message: "This field must be at least 3 characters long",
					},
				},
			},
		},
		{
			name:     "validation error - all fields empty",
			userID:   1,
			courseID: "57",
			requestBody: updateRequest{
				Title:       "",
				Description: "",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "title",
						Message: "This field is required",
					},
					{
						Field:   "description",
						Message: "This field is required",
					},
				},
			},
		},
		{
			name:     "invalid course id - not a number",
			userID:   1,
			courseID: "ivan",
			requestBody: updateRequest{
				Title:       "Course #1",
				Description: "The best course you can choose!",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - zero",
			userID:   1,
			courseID: "0",
			requestBody: updateRequest{
				Title:       "Course #1",
				Description: "The best course you can choose!",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "invalid course id - negative",
			userID:   1,
			courseID: "-12",
			requestBody: updateRequest{
				Title:       "Course #1",
				Description: "The best course you can choose!",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "course id must be a positive integer",
			},
		},
		{
			name:     "missing user id in context",
			noCtxID:  true,
			courseID: "30",
			requestBody: updateRequest{
				Title:       "New course",
				Description: "learn more",
			},
			mockFn:   func(m *mockCourseService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:     "course not found",
			userID:   1,
			courseID: "122226",
			requestBody: updateRequest{
				Title:       "New course",
				Description: "learn more",
			},
			mockFn: func(m *mockCourseService) {
				in := entity.Course{
					ID:          122226,
					TeacherID:   1,
					Title:       "New course",
					Description: "learn more",
				}
				m.On("Update", mock.Anything, in).
					Return(entity.Course{}, apperror.CourseNotFoundID(122226)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "course with id 122226 was not found",
			},
		},
		{
			name:     "access denied - not course owner",
			userID:   2,
			courseID: "7",
			requestBody: updateRequest{
				Title:       "My course",
				Description: "I'm pretty sure you will like it.",
			},
			mockFn: func(m *mockCourseService) {
				in := entity.Course{
					ID:          7,
					TeacherID:   2,
					Title:       "My course",
					Description: "I'm pretty sure you will like it.",
				}
				m.On("Update", mock.Anything, in).
					Return(entity.Course{}, apperror.CourseAccessDenied).Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: httpresponse.ErrorResponse{
				Message: "only author can manage this course",
			},
		},
		{
			name:     "service returns internal error",
			userID:   1,
			courseID: "8",
			requestBody: updateRequest{
				Title:       "Something",
				Description: "something",
			},
			mockFn: func(m *mockCourseService) {
				in := entity.Course{
					ID:          8,
					TeacherID:   1,
					Title:       "Something",
					Description: "something",
				}
				m.On("Update", mock.Anything, in).
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

			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodPut, "/courses/"+tt.courseID, bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")
			r.SetPathValue("id", tt.courseID)

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.update)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantJSON), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
