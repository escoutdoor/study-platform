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

func TestCreate(t *testing.T) {
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
			userID: 1,
			body: createRequest{
				Department: "IT department",
			},
			mockFn: func(m *mockTeacherService) {
				in := entity.Teacher{
					UserID:     1,
					Department: "IT department",
				}
				m.On("Create", mock.Anything, in).Return(nil).Once()
			},
			wantCode: http.StatusCreated,
			wantBody: map[string]string{"message": "teacher profile created successfully"},
		},
		{
			name:     "invalid json",
			userID:   1,
			body:     "hahahaha not json, i will break everything",
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "invalid request body format",
			},
		},
		{
			name:   "validation error - empty department",
			userID: 1,
			body: createRequest{
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
			name:   "validation error - department length < 2 symbols",
			userID: 1,
			body: createRequest{
				Department: "T",
			},
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "request validation failed",
				Details: []httpresponse.ErrorDetail{
					{
						Field:   "department",
						Message: "This field must be at least 2 characters long",
					},
				},
			},
		},
		{
			name: "no user id in ctx",
			body: createRequest{
				Department: "Something",
			},
			noCtxID:  true,
			mockFn:   func(m *mockTeacherService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
		{
			name:   "teacher already exists",
			userID: 1,
			body: createRequest{
				Department: "Computer Science",
			},
			mockFn: func(m *mockTeacherService) {
				in := entity.Teacher{UserID: 1, Department: "Computer Science"}
				m.On("Create", mock.Anything, in).
					Return(apperror.TeacherAlreadyExists).Once()
			},
			wantCode: http.StatusConflict,
			wantBody: httpresponse.ErrorResponse{
				Message: "you are already a teacher",
			},
		},
		{
			name:   "service returns internal error",
			userID: 1,
			body: createRequest{
				Department: "Computer Science",
			},
			mockFn: func(m *mockTeacherService) {
				in := entity.Teacher{UserID: 1, Department: "Computer Science"}
				m.On("Create", mock.Anything, in).
					Return(errors.New("this field doesn't exist in your database")).Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpresponse.ErrorResponse{
				Message: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		svc := new(mockTeacherService)
		cv := validator.New()
		h := &handler{service: svc, cv: cv}

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

		r := httptest.NewRequest(http.MethodPost, "/teachers", bytes.NewReader(requestBodyBytes))
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()

		if !tt.noCtxID {
			ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, tt.userID)
			r = r.WithContext(ctx)
		}

		errhandler.ErrorHandler(h.create).ServeHTTP(w, r)

		assert.Equal(t, tt.wantCode, w.Code)

		wantBodyJson, err := json.Marshal(tt.wantBody)
		require.NoError(t, err)
		require.JSONEq(t, string(wantBodyJson), string(wantBodyJson))

		svc.AssertExpectations(t)
	}
}
