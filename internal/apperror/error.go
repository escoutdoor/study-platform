package apperror

import (
	"errors"
	"fmt"

	"github.com/escoutdoor/study-platform/internal/apperror/code"
)

var (
	ErrInvalidJSON = newError(code.InvalidJson, "invalid request body format")

	CourseAccessDenied  = newError(code.PermissionDenied, "only author can manage this course")
	StudentAccessDenied = newError(code.PermissionDenied, "you can only manage your own student profile")
	TeacherAccessDenied = newError(code.PermissionDenied, "you can only manage your own teacher profile")

	StudentAlreadyEnrolled = newError(code.StudentAlreadyEnrolled, "you are already enrolled in this course")
	StudentNotEnrolled     = newError(code.StudentNotEnrolled, "you are not enrolled in this course")
	TeacherAlreadyExists   = newError(code.TeacherAlreadyExists, "you are already a teacher")

	ErrJwtTokenExpired      = newError(code.JwtTokenExpired, "jwt token is already expired")
	ErrInvalidJwtToken      = newError(code.InvalidJwtToken, "invalid jwt token")
	ErrIncorrectCredentials = newError(code.IncorrectCredentials, "incorrect credentials")
)

type Error struct {
	Code code.Code
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func newError(code code.Code, err string) *Error {
	return &Error{
		Code: code,
		Err:  errors.New(err),
	}
}

func UserEmailAlreadyExists(email string) *Error {
	msg := fmt.Sprintf("user with email %q is already exists", email)
	return newError(code.EmailAlreadyExists, msg)
}

func StudentEmailAlreadyExists(email string) *Error {
	msg := fmt.Sprintf("student with email %q is already exists", email)
	return newError(code.EmailAlreadyExists, msg)
}

func TeacherEmailAlreadyExists(email string) *Error {
	msg := fmt.Sprintf("teacher with email %q is already exists", email)
	return newError(code.EmailAlreadyExists, msg)
}

func ValidationFailed(msg string) *Error {
	return newError(code.ValidationFailed, msg)
}

func StudentNotFoundID(studentID int) *Error {
	msg := fmt.Sprintf("student with id %d was not found", studentID)
	return newError(code.StudentNotFound, msg)
}

func StudentNotFoundEmail(email string) *Error {
	msg := fmt.Sprintf("student with email %q was not found", email)
	return newError(code.StudentNotFound, msg)
}

func UserNotFoundEmail(email string) *Error {
	msg := fmt.Sprintf("user with email %q was not found", email)
	return newError(code.UserNotFound, msg)
}

func TeacherNotFoundID(teacherID int) *Error {
	msg := fmt.Sprintf("teacher with id %d was not found", teacherID)
	return newError(code.TeacherNotFound, msg)
}

func TeacherNotFoundEmail(email string) *Error {
	msg := fmt.Sprintf("teacher with email %q was not found", email)
	return newError(code.TeacherNotFound, msg)
}

func CourseNotFoundID(courseID int) *Error {
	msg := fmt.Sprintf("course with id %d was not found", courseID)
	return newError(code.CourseNotFound, msg)
}

func UserNotFoundID(userID int) *Error {
	msg := fmt.Sprintf("user with id %d was not found", userID)
	return newError(code.UserNotFound, msg)
}
