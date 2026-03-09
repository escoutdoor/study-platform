package student

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
	studentService    studentService
	enrollmentService enrollmentService
	cv                *validator.CustomValidator
}

func RegisterHandlers(
	mux *http.ServeMux,
	studentService studentService,
	enrollmentService enrollmentService,
	cv *validator.CustomValidator,
	authMiddleware func(errhandler.HandlerFunc) errhandler.HandlerFunc,
) {
	h := &handler{
		studentService:    studentService,
		enrollmentService: enrollmentService,
		cv:                cv,
	}

	routes := map[string]errhandler.HandlerFunc{
		"GET /students":      h.list,
		"GET /students/{id}": h.get,

		"PUT /students/me": authMiddleware(middleware.RequireRole(token.RoleStudent)(h.update)),

		// optional
		"POST /students/me/courses/{courseId}": authMiddleware(middleware.RequireRole(token.RoleStudent)(h.enroll)),

		"DELETE /students/me/courses/{courseId}": authMiddleware(middleware.RequireRole(token.RoleStudent)(h.unenroll)),
	}
	for p, h := range routes {
		mux.Handle(p, errhandler.ErrorHandler(h))
	}
}

type studentService interface {
	List(ctx context.Context) ([]entity.Student, error)
	Get(ctx context.Context, userID int) (entity.Student, error)

	Update(ctx context.Context, in entity.Student) error
}

type enrollmentService interface {
	Enroll(ctx context.Context, userID, courseID int) error
	Unenroll(ctx context.Context, userID, courseID int) error
}

type studentResponse struct {
	UserID int `json:"userId"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email string `json:"email"`
}

func studentToResponse(student entity.Student) studentResponse {
	return studentResponse{
		UserID: student.UserID,

		FirstName: student.FirstName,
		LastName:  student.LastName,

		Email: student.Email,
	}
}

func studentListToResponse(students []entity.Student) []studentResponse {
	list := make([]studentResponse, 0, len(students))
	for _, s := range students {
		list = append(list, studentToResponse(s))
	}

	return list
}

func validateUserID(r *http.Request) (int, error) {
	idParamStr := r.PathValue(idParam)
	userID, ok := isValidID(idParamStr)
	if !ok {
		return 0, apperror.ValidationFailed("student id must be a positive integer")
	}

	return userID, nil
}

func validateCourseID(r *http.Request) (int, error) {
	idParamStr := r.PathValue("courseId")
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
