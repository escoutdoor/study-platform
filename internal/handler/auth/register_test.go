package auth

import (
	"bytes"
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

func TestRegister(t *testing.T) {

	tests := []struct {
		name string

		body   any
		mockFn func(m *mockAuthService)

		wantCode int
		wantBody any
	}{
		{
			name: "success",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "hello@gmail.com",
				Password:  "Hello123!",
			},
			mockFn: func(m *mockAuthService) {
				m.On("Register", mock.Anything, entity.User{
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "hello@gmail.com",
					Password:  "Hello123!",
				}).Return(entity.Tokens{
					AccessToken:  "valid access token",
					RefreshToken: "valid refresh token",
				}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: registerResponse{
				Tokens: authResponse{
					AccessToken:  "valid access token",
					RefreshToken: "valid refresh token",
				},
			},
		},
		{
			name:     "invalid json",
			body:     "not a json{12121{{",
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name: "validation error - empty first name",
			body: registerRequest{
				FirstName: "",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "password123",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "firstName",
						Message: "This field is required",
					},
				},
			},
		},
		{
			name: "validation error - first name too long",
			body: registerRequest{
				FirstName: "IvanIvanIvanIvanIvanI",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "password123",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "firstName",
						Message: "This field must be at most 20 characters long",
					},
				},
			},
		},
		{
			name: "validation error - empty last name",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "",
				Email:     "ivan@example.com",
				Password:  "password123",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "lastName", Message: "This field is required"},
				},
			},
		},
		{
			name: "validation error - invalid email",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "not-an-email",
				Password:  "password123",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "email", Message: "This field must be a valid email address"},
				},
			},
		},
		{
			name: "validation error - empty email",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "",
				Password:  "password123",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "email", Message: "This field is required"},
				},
			},
		},
		{
			name: "validation error - empty password",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "password", Message: "This field is required"},
				},
			},
		},
		{
			name: "validation error - password too short",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "short",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "password", Message: "This field must be at least 8 characters long"},
				},
			},
		},
		{
			name: "validation error - password too long",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "thispasswordiswaytoolongthispasswordiswaytoolongthispasswordiswaytoolong",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "password", Message: "This field must be at most 40 characters long"},
				},
			},
		},
		{
			name: "email already exists",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "password123",
			},
			mockFn: func(m *mockAuthService) {
				in := entity.User{
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
					Password:  "password123",
				}
				m.On("Register", mock.Anything, in).
					Return(entity.Tokens{}, apperror.UserEmailAlreadyExists("ivan@example.com")).Once()
			},
			wantCode: http.StatusConflict,
			wantBody: httpresponse.ErrorResponse{
				Message: `user with email "ivan@example.com" is already exists`,
			},
		},
		{
			name: "service returns internal error",
			body: registerRequest{
				FirstName: "Ivan",
				LastName:  "Popov",
				Email:     "ivan@example.com",
				Password:  "password123",
			},
			mockFn: func(m *mockAuthService) {
				in := entity.User{
					FirstName: "Ivan",
					LastName:  "Popov",
					Email:     "ivan@example.com",
					Password:  "password123",
				}
				m.On("Register", mock.Anything, in).
					Return(entity.Tokens{}, errors.New("db connection lost")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockAuthService)
			cv := validator.New()
			h := &handler{authService: svc, cv: cv}
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

			r := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.register)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJson, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantJson), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
