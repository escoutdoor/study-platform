package user

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
	tests := []struct {
		name string

		userID  int
		noCtxID bool
		mockFn  func(m *mockUserService)

		wantCode int
		wantBody *httpresponse.ErrorResponse
	}{
		{
			name:   "success",
			userID: 12,
			mockFn: func(m *mockUserService) {
				m.On("Delete", mock.Anything, 12).Return(nil).Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:     "missing user id in context",
			noCtxID:  true,
			mockFn:   func(m *mockUserService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: &httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "user not found",
			userID: 12,
			mockFn: func(m *mockUserService) {
				m.On("Delete", mock.Anything, 12).
					Return(apperror.UserNotFoundID(12)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: &httpresponse.ErrorResponse{
				Message: "user with id 12 was not found",
			},
		},
		{
			name:   "service returns internal error",
			userID: 1,
			mockFn: func(m *mockUserService) {
				m.On("Delete", mock.Anything, 1).
					Return(errors.New("something bad happened")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockUserService)
			cv := validator.New()
			h := &handler{userService: svc, cv: cv}

			tt.mockFn(svc)

			r := httptest.NewRequest(http.MethodDelete, "/users/me", nil)
			w := httptest.NewRecorder()

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}

			errhandler.ErrorHandler(h.delete).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantBody != nil {
				wantBodyJson, err := json.Marshal(tt.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantBodyJson), w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}

			svc.AssertExpectations(t)
		})
	}
}
