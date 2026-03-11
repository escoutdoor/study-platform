package teacher

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"github.com/escoutdoor/study-platform/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tests := []struct {
		name string

		mockFn func(m *mockTeacherService)

		wantCode int
		wantBody any
	}{
		{
			name: "success - multiple teachers",
			mockFn: func(m *mockTeacherService) {
				m.On("List", mock.Anything).Return([]entity.Teacher{
					{
						UserID:     1,
						FirstName:  "Ivan",
						LastName:   "Popov",
						Department: "Computer Science",
						Email:      "ivan@example.com",
					},
					{
						UserID:     2,
						FirstName:  "Masha",
						LastName:   "Dnipro",
						Department: "Mathematics",
						Email:      "masha@example.com",
					},
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Teachers: []teacherResponse{
					{UserID: 1, FirstName: "Ivan", LastName: "Popov", Department: "Computer Science", Email: "ivan@example.com"},
					{UserID: 2, FirstName: "Masha", LastName: "Dnipro", Department: "Mathematics", Email: "masha@example.com"},
				},
			},
		},
		{
			name: "success - empty list",
			mockFn: func(m *mockTeacherService) {
				m.On("List", mock.Anything).Return([]entity.Teacher{}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listResponse{
				Teachers: []teacherResponse{},
			},
		},
		{
			name: "service returns internal error",
			mockFn: func(m *mockTeacherService) {
				m.On("List", mock.Anything).
					Return([]entity.Teacher(nil), errors.New("bad connection")).Once()
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

			r := httptest.NewRequest(http.MethodGet, "/teachers", nil)
			w := httptest.NewRecorder()

			errhandler.ErrorHandler(h.list).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), w.Body.String())

			svc.AssertExpectations(t)
		})
	}
}
