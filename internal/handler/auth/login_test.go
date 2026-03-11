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

func TestLogin(t *testing.T) {

	tests := []struct {
		name string

		body   any
		mockFn func(m *mockAuthService)

		wantCode int
		wantBody any
	}{
		{
			name: "success",
			body: loginRequest{
				Email:    "vanap387@gmail.com",
				Password: "myproject123!",
			},
			mockFn: func(m *mockAuthService) {
				in := entity.User{
					Email:    "vanap387@gmail.com",
					Password: "myproject123!",
				}
				m.On("Login", mock.Anything, in).Return(entity.Tokens{
					AccessToken:  "valid access token",
					RefreshToken: "valid refresh token",
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: loginResponse{
				Tokens: authResponse{
					AccessToken:  "valid access token",
					RefreshToken: "valid refresh token",
				},
			},
		},
		{
			name:     "invalid json",
			body:     "i like it{{{",
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name: "validation error - empty email",
			body: loginRequest{
				Email:    "",
				Password: "password123",
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
			name: "validation error - invalid email format",
			body: loginRequest{
				Email:    "not-an-email",
				Password: "password123",
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
			name: "validation error - empty password",
			body: loginRequest{
				Email:    "ivan@example.com",
				Password: "",
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
			body: loginRequest{
				Email:    "ivan@example.com",
				Password: "short",
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
			body: loginRequest{
				Email:    "ivan@example.com",
				Password: "thispasswordiswaytoolongthispasswordiswaytoolongthispasswordiswaytoolonghello",
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
			name: "incorrect credentials",
			body: loginRequest{
				Email:    "ivan@example.com",
				Password: "wrongpassword1",
			},
			mockFn: func(m *mockAuthService) {
				in := entity.User{
					Email:    "ivan@example.com",
					Password: "wrongpassword1",
				}
				m.On("Login", mock.Anything, in).
					Return(entity.Tokens{}, apperror.ErrIncorrectCredentials).Once()
			},
			wantCode: http.StatusUnauthorized,
			wantBody: httpresponse.ErrorResponse{
				Message: "incorrect credentials",
			},
		},
		{
			name: "service returns internal error",
			body: loginRequest{
				Email:    "ivan@example.com",
				Password: "password123hello!!!H",
			},
			mockFn: func(m *mockAuthService) {
				in := entity.User{
					Email:    "ivan@example.com",
					Password: "password123hello!!!H",
				}
				m.On("Login", mock.Anything, in).
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
			h := &handler{
				authService: svc,
				cv:          cv,
			}
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

			r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.login)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJson, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantJson), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
