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

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name string

		body   any
		mockFn func(m *mockAuthService)

		wantCode int
		wantBody any
	}{
		{
			name: "success",
			body: refreshTokenRequest{
				RefreshToken: "valid refresh token",
			},
			mockFn: func(m *mockAuthService) {
				m.
					On(
						"RefreshToken",
						mock.Anything,
						"valid refresh token",
					).
					Return(
						entity.Tokens{
							AccessToken:  "cool access token",
							RefreshToken: "the best refresh token",
						},
						nil,
					)
			},
			wantCode: http.StatusOK,
			wantBody: refreshTokenResponse{
				Tokens: authResponse{
					AccessToken:  "cool access token",
					RefreshToken: "the best refresh token",
				},
			},
		},
		{
			name:     "invalid json",
			body:     "bad json so bad{{{",
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name: "validation error - empty refresh token",
			body: refreshTokenRequest{
				RefreshToken: "",
			},
			mockFn:   func(m *mockAuthService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "refreshToken", Message: "This field is required"},
				},
			},
		},
		{
			name: "expired refresh token",
			body: refreshTokenRequest{
				RefreshToken: "expired-token",
			},
			mockFn: func(m *mockAuthService) {
				m.On("RefreshToken", mock.Anything, "expired-token").
					Return(entity.Tokens{}, apperror.ErrJwtTokenExpired).Once()
			},
			wantCode: http.StatusUnauthorized,
			wantBody: httpresponse.ErrorResponse{
				Message: "jwt token is already expired",
			},
		},
		{
			name: "invalid refresh token",
			body: refreshTokenRequest{
				RefreshToken: "invalid-token",
			},
			mockFn: func(m *mockAuthService) {
				m.On("RefreshToken", mock.Anything, "invalid-token").
					Return(entity.Tokens{}, apperror.ErrInvalidJwtToken).Once()
			},
			wantCode: http.StatusUnauthorized,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid jwt token",
			},
		},
		{
			name: "user not found",
			body: refreshTokenRequest{
				RefreshToken: "token-for-deleted-user",
			},
			mockFn: func(m *mockAuthService) {
				m.On("RefreshToken", mock.Anything, "token-for-deleted-user").
					Return(entity.Tokens{}, apperror.UserNotFoundID(123)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "user with id 123 was not found",
			},
		},
		{
			name: "service returns internal error",
			body: refreshTokenRequest{
				RefreshToken: "some-token",
			},
			mockFn: func(m *mockAuthService) {
				m.On("RefreshToken", mock.Anything, "some-token").
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
			h := &handler{authService: svc, cv: validator.New()}
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

			r := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			httpHandler := errhandler.ErrorHandler(h.refreshToken)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantBodyJson, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantBodyJson), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
