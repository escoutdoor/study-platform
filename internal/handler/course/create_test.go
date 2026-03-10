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

func TestCreate(t *testing.T) {
	type mockBehavior func(m *mockCourseService, in entity.Course)

	tests := []struct {
		name string

		userID  int
		body    any
		mockFn  mockBehavior
		noCtxID bool

		wantCode int
		wantBody any
	}{
		{
			name:   "success",
			userID: 1,
			body: createRequest{
				Title:       "Course #1",
				Description: "My first course",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {
				m.On("Create", mock.Anything, in).Return(10, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: createResponse{CourseID: 10},
		},
		{
			name:   "invalid json",
			userID: 1,
			body:   "not a json{{{",
			mockFn: func(m *mockCourseService, in entity.Course) {},

			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name:   "validation error - empty title",
			userID: 1,
			body: createRequest{
				Title:       "",
				Description: "some description",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {},

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
			name:   "validation error - title too short",
			userID: 1,
			body: createRequest{
				Title:       "IP",
				Description: "some description",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {},

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
			name:   "validation error - empty description",
			userID: 1,
			body: createRequest{
				Title:       "My favorite course",
				Description: "",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {},

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
			name:   "validation error - all fields empty",
			userID: 1,
			body: createRequest{
				Title:       "",
				Description: "",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {},

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
			name:    "missing user id in context",
			noCtxID: true,
			body: createRequest{
				Title:       "My first course",
				Description: "the best course ever",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {},

			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "service returns internal error",
			userID: 1,
			body: createRequest{
				Title:       "Title of the course",
				Description: "some description",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {
				m.On("Create", mock.Anything, in).Return(0, errors.New("db connection lost"))
			},

			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "service returns app error",
			userID: 1,
			body: createRequest{
				Title:       "How to cook?",
				Description: "I like it",
			},
			mockFn: func(m *mockCourseService, in entity.Course) {
				m.On("Create", mock.Anything, in).
					Return(0, apperror.UserNotFoundID(1))
			},

			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "user with id 1 was not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockCourseService)
			cv := validator.New()
			h := &handler{service: svc, cv: cv}

			in := entity.Course{}
			if req, ok := tt.body.(createRequest); ok && !tt.noCtxID {
				in = createRequestToCourse(&req, tt.userID)
			}
			tt.mockFn(svc, in)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.create)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantBody != nil {
				wantJSON, err := json.Marshal(tt.wantBody)
				require.NoError(t, err)

				assert.JSONEq(t, string(wantJSON), w.Body.String())
			}

			svc.AssertExpectations(t)
		})
	}
}
