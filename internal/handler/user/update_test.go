package user

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
		noCtxID bool
		body    any
		mockFn  func(m *mockUserService)

		wantCode int
		wantBody any
	}{
		{
			name:   "success",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn: func(m *mockUserService) {
				in := entity.User{
					ID:        1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
					Password:  "newpass123",
				}
				m.On("Update", mock.Anything, in).Return(entity.User{
					ID:        1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: updateResponse{
				User: userResponse{
					ID:        1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
				},
			},
		},
		{
			name:     "invalid json",
			userID:   1,
			body:     "not a json",
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name:   "validation error - empty first name",
			userID: 1,
			body: updateRequest{
				FirstName: "",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "firstName", Message: "This field is required"},
				},
			},
		},
		{
			name:   "validation error - first name too long",
			userID: 1,
			body: updateRequest{
				FirstName: "IvanIvanIvanIvanIvanI",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "firstName", Message: "This field must be at most 20 characters long"},
				},
			},
		},
		{
			name:   "validation error - empty last name",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "lastName", Message: "This field is required"},
				},
			},
		},
		{
			name:   "validation error - last name too long",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "PopovPopovPopovPopovP",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "lastName", Message: "This field must be at most 20 characters long"},
				},
			},
		},
		{
			name:   "validation error - empty email",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "email", Message: "This field is required"},
				},
			},
		},
		{
			name:   "validation error - invalid email format",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "not-an-email",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "email", Message: "This field must be a valid email address"},
				},
			},
		},
		{
			name:   "validation error - empty password",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "password", Message: "This field is required"},
				},
			},
		},
		{
			name:   "validation error - password too short",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "short",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "password", Message: "This field must be at least 8 characters long"},
				},
			},
		},
		{
			name:   "validation error - password too long",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "thispasswordiswaytoolongforthelimitof40chars",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "password", Message: "This field must be at most 40 characters long"},
				},
			},
		},
		{
			name:    "missing user id in context",
			noCtxID: true,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "user not found",
			userID: 88,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn: func(m *mockUserService) {
				in := entity.User{
					ID:        88,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
					Password:  "newpass123",
				}
				m.On("Update", mock.Anything, in).
					Return(entity.User{}, apperror.UserNotFoundID(88)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "user with id 88 was not found",
			},
		},
		{
			name:   "email already exists",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "danya@example.com",
				Password:  "newpass123",
			},
			mockFn: func(m *mockUserService) {
				in := entity.User{
					ID:        1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "danya@example.com",
					Password:  "newpass123",
				}
				m.On("Update", mock.Anything, in).
					Return(entity.User{}, apperror.UserEmailAlreadyExists("danya@example.com")).Once()
			},
			wantCode: http.StatusConflict,
			wantBody: httpresponse.ErrorResponse{
				Message: `user with email "danya@example.com" is already exists`,
			},
		},
		{
			name:   "service returns internal error",
			userID: 1,
			body: updateRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "newpass123",
			},
			mockFn: func(m *mockUserService) {
				in := entity.User{
					ID:        1,
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
					Password:  "newpass123",
				}
				m.On("Update", mock.Anything, in).
					Return(entity.User{}, errors.New("something happened and we don't know what")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockUserService)
			h := &handler{userService: svc, cv: validator.New()}

			tt.mockFn(svc)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodPut, "/users/me", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.updateMe).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
