package teacher

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
		mockFn  func(m *mockTeacherService)

		wantCode int
		wantBody any
	}{
		{
			name:   "success",
			userID: 40,
			body: updateRequest{
				Department: "Math",
			},
			mockFn: func(m *mockTeacherService) {
				in := entity.Teacher{
					UserID:     40,
					Department: "Math",
				}
				m.On("Update", mock.Anything, in).Return(
					entity.Teacher{
						UserID:     40,
						FirstName:  "Ivan",
						LastName:   "Popov",
						Department: "Math",
						Email:      "vanek@example.com",
					},
					nil,
				).Once()
			},
			wantCode: http.StatusOK,
			wantBody: updateResponse{
				teacherResponse{
					UserID:     40,
					FirstName:  "Ivan",
					LastName:   "Popov",
					Department: "Math",
					Email:      "vanek@example.com",
				},
			},
		},
		{
			name:     "invalid json",
			userID:   1,
			body:     "not a json{{{",
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name:   "validation error - empty department",
			userID: 1,
			body: updateRequest{
				Department: "",
			},
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "department", Message: "This field is required"},
				},
			},
		},
		{
			name:   "validation error - department too short",
			userID: 1,
			body: updateRequest{
				Department: "T",
			},
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{Field: "department", Message: "This field must be at least 2 characters long"},
				},
			},
		},
		{
			name:    "missing user id in context",
			noCtxID: true,
			body: updateRequest{
				Department: "Mathematics",
			},
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "teacher not found",
			userID: 8,
			body: updateRequest{
				Department: "Mathematics",
			},
			mockFn: func(m *mockTeacherService) {
				in := entity.Teacher{UserID: 8, Department: "Mathematics"}
				m.On("Update", mock.Anything, in).
					Return(entity.Teacher{}, apperror.TeacherNotFoundID(8)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "teacher with id 8 was not found",
			},
		},
		{
			name:   "service returns internal error",
			userID: 1,
			body: updateRequest{
				Department: "Mathematics",
			},
			mockFn: func(m *mockTeacherService) {
				in := entity.Teacher{UserID: 1, Department: "Mathematics"}
				m.On("Update", mock.Anything, in).
					Return(entity.Teacher{}, errors.New("so bad error")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockTeacherService)
			cv := validator.New()
			h := &handler{
				service: svc,
				cv:      cv,
			}
			tt.mockFn(svc)

			var requestBodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				requestBodyBytes = []byte(v)
			default:
				var err error
				requestBodyBytes, err = json.Marshal(v)
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodPut, "/teachers/me", bytes.NewReader(requestBodyBytes))
			r.Header.Add("Content-Type", "application/json")

			if !tt.noCtxID {
				ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.update).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantBodyJson, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			require.JSONEq(t, string(wantBodyJson), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
