package teacher

import (
	"context"
	"net/http"
	"strconv"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/middleware"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/internal/util/token"
	"github.com/escoutdoor/study-platform/pkg/validator"
)

const (
	idParam = "id"
)

type handler struct {
	service teacherService
	cv      *validator.CustomValidator
}

func RegisterHandlers(
	mux *http.ServeMux,
	teacherService teacherService,
	cv *validator.CustomValidator,
	authMiddleware func(errhandler.HandlerFunc) errhandler.HandlerFunc,
) {
	h := &handler{service: teacherService, cv: cv}

	routes := map[string]errhandler.HandlerFunc{
		"GET /teachers":      h.list,
		"GET /teachers/{id}": h.get,

		"POST /teachers": authMiddleware(h.create),

		"PUT /teachers/me": authMiddleware(middleware.RequireRole(token.RoleTeacher)(h.update)),
	}
	for p, h := range routes {
		mux.Handle(p, errhandler.ErrorHandler(h))
	}
}

type teacherService interface {
	List(ctx context.Context) ([]entity.Teacher, error)
	Get(ctx context.Context, userID int) (entity.Teacher, error)

	Create(ctx context.Context, in entity.Teacher) error
	Update(ctx context.Context, in entity.Teacher) (entity.Teacher, error)
}

type teacherResponse struct {
	UserID int `json:"userId"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Department string `json:"department"`

	Email string `json:"email"`
}

func teacherToResponse(teacher entity.Teacher) teacherResponse {
	return teacherResponse{
		UserID: teacher.UserID,

		FirstName: teacher.FirstName,
		LastName:  teacher.LastName,

		Department: teacher.Department,

		Email: teacher.Email,
	}
}

func teacherListToResponse(teachers []entity.Teacher) []teacherResponse {
	list := make([]teacherResponse, 0, len(teachers))
	for _, s := range teachers {
		list = append(list, teacherToResponse(s))
	}

	return list
}

func validateUserID(r *http.Request) (int, error) {
	idParamStr := r.PathValue(idParam)
	userID, ok := isValidID(idParamStr)
	if !ok {
		return 0, apperror.ValidationFailed("teacher id must be a positive integer")
	}

	return userID, nil
}

func isValidID(param string) (int, bool) {
	id, err := strconv.Atoi(param)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
