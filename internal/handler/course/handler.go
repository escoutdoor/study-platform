package course

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
	service courseService
	cv      *validator.CustomValidator
}

func RegisterHandlers(
	mux *http.ServeMux,
	courseService courseService,
	cv *validator.CustomValidator,
	authMiddleware func(errhandler.HandlerFunc) errhandler.HandlerFunc,
) {
	h := &handler{service: courseService, cv: cv}

	routes := map[string]errhandler.HandlerFunc{
		"GET /courses":      h.list,
		"GET /courses/{id}": h.get,

		"POST /courses": authMiddleware(middleware.RequireRole(token.RoleTeacher)(h.create)),

		"PUT /courses/{id}": authMiddleware(middleware.RequireRole(token.RoleTeacher)(h.update)),

		"DELETE /courses/{id}": authMiddleware(middleware.RequireRole(token.RoleTeacher)(h.delete)),
	}
	for p, h := range routes {
		mux.Handle(p, errhandler.ErrorHandler(h))
	}
}

type courseService interface {
	List(ctx context.Context) ([]entity.Course, error)
	Get(ctx context.Context, courseID int) (entity.Course, error)

	Create(ctx context.Context, in entity.Course) (int, error)
	Update(ctx context.Context, in entity.Course) (entity.Course, error)
	Delete(ctx context.Context, courseID int, teacherID int) error
}

type courseResponse struct {
	ID        int `json:"id"`
	TeacherID int `json:"teacherId"`

	Title       string `json:"title"`
	Description string `json:"description"`
}

func courseToResponse(course entity.Course) courseResponse {
	return courseResponse{
		ID:        course.ID,
		TeacherID: course.TeacherID,

		Title:       course.Title,
		Description: course.Description,
	}
}

func courseListToResponse(courses []entity.Course) []courseResponse {
	list := make([]courseResponse, 0, len(courses))
	for _, s := range courses {
		list = append(list, courseToResponse(s))
	}

	return list
}

func validateCourseID(r *http.Request) (int, error) {
	idParamStr := r.PathValue(idParam)
	courseID, ok := isValidID(idParamStr)
	if !ok {
		return 0, apperror.ValidationFailed("course id must be a positive integer")
	}

	return courseID, nil
}

func isValidID(param string) (int, bool) {
	id, err := strconv.Atoi(param)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
