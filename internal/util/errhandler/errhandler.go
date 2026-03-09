package errhandler

import (
	"errors"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/apperror/code"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"github.com/escoutdoor/study-platform/pkg/logger"
	"github.com/escoutdoor/study-platform/pkg/validator"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		ctx := r.Context()

		var validationErr *validator.ValidationError
		if errors.As(err, &validationErr) {
			if r.Method == http.MethodHead {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			details := make([]httpresponse.ErrorDetail, 0, len(validationErr.Errors))
			for _, item := range validationErr.Errors {
				details = append(details, httpresponse.ErrorDetail{
					Field:   item.Field,
					Message: item.Message,
				})
			}

			httpresponse.JSON(w, http.StatusBadRequest, httpresponse.ErrorResponse{
				Message: validationErr.Error(),
				Details: details,
			})
			return
		}

		respCode := http.StatusInternalServerError
		resp := httpresponse.ErrorResponse{
			Message: "internal server error",
		}

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case code.EmailAlreadyExists,
				code.StudentAlreadyEnrolled,
				code.StudentNotEnrolled,
				code.TeacherAlreadyExists:
				respCode = http.StatusConflict

			case code.PermissionDenied:
				respCode = http.StatusForbidden

			case code.StudentNotFound,
				code.TeacherNotFound,
				code.CourseNotFound,
				code.UserNotFound:
				respCode = http.StatusNotFound

			case code.InvalidJson, code.ValidationFailed:
				respCode = http.StatusBadRequest

			case code.IncorrectCredentials,
				code.JwtTokenExpired,
				code.InvalidJwtToken:
				respCode = http.StatusUnauthorized

			default:
				respCode = http.StatusInternalServerError
			}

			resp = httpresponse.ErrorResponse{
				Message: appErr.Error(),
			}
		}

		if respCode == http.StatusInternalServerError {
			logger.ErrorKV(ctx, "internal server error", "err", err.Error())
		}

		if r.Method == http.MethodHead {
			w.WriteHeader(respCode)
			return
		}

		httpresponse.JSON(w, respCode, resp)
	})
}
