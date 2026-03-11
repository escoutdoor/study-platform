package teacher

import (
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

func TestGet(t *testing.T) {
	tests := []struct {
		name string

		teacherID string
		mockFn    func(m *mockTeacherService)

		wantCode int
		wantBody any
	}{
		{
			name:      "success",
			teacherID: "1",
			mockFn: func(m *mockTeacherService) {
				m.On("Get", mock.Anything, 1).Return(
					entity.Teacher{
						UserID:     1,
						FirstName:  "Hello",
						LastName:   "Super Hello",
						Department: "IT Department",
						Email:      "example@example.com",
					},
					nil,
				).Once()
			},
			wantCode: http.StatusOK,
			wantBody: getResponse{
				Teacher: teacherResponse{
					UserID:     1,
					FirstName:  "Hello",
					LastName:   "Super Hello",
					Department: "IT Department",
					Email:      "example@example.com",
				},
			},
		},
		{
			name:      "invalid teacher id - not a number",
			teacherID: "popov",
			mockFn:    func(m *mockTeacherService) {},
			wantCode:  http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "teacher id must be a positive integer",
			},
		},
		{
			name:      "invalid teacher id - zero",
			teacherID: "0",
			mockFn:    func(m *mockTeacherService) {},
			wantCode:  http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "teacher id must be a positive integer",
			},
		},
		{
			name:      "invalid teacher id - negative",
			teacherID: "-363",
			mockFn:    func(m *mockTeacherService) {},
			wantCode:  http.StatusBadRequest,
			wantBody: httpresponse.ErrorResponse{
				Message: "teacher id must be a positive integer",
			},
		},
		{
			name:      "teacher not found",
			teacherID: "9833",
			mockFn: func(m *mockTeacherService) {
				m.On("Get", mock.Anything, 9833).
					Return(entity.Teacher{}, apperror.TeacherNotFoundID(9833)).Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: httpresponse.ErrorResponse{
				Message: "teacher with id 9833 was not found",
			},
		},
		{
			name:      "service returns internal error",
			teacherID: "1",
			mockFn: func(m *mockTeacherService) {
				m.On("Get", mock.Anything, 1).
					Return(entity.Teacher{}, errors.New("this table doesn't exist")).Once()
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
			h := &handler{service: svc, cv: validator.New()}

			tt.mockFn(svc)

			r := httptest.NewRequest(http.MethodGet, "/teachers/"+tt.teacherID, nil)
			r.SetPathValue("id", tt.teacherID)

			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.get).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
