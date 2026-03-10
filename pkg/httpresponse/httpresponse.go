package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/pkg/logger"
)

type ErrorResponse struct {
	Message string        `json:"message" example:"request validation failed"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty" example:"email"`
	Message string `json:"message" example:"This field is required"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	if v == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.DebugKV(context.Background(), "encode json response", "error", err)
	}
}

func Created(w http.ResponseWriter, v any) {
	JSON(w, http.StatusCreated, v)
}

func OK(w http.ResponseWriter, v any) {
	JSON(w, http.StatusOK, v)
}

func NoContent(w http.ResponseWriter) {
	JSON(w, http.StatusNoContent, nil)
}

//

func Error(w http.ResponseWriter, status int, err string) {
	JSON(w, status, ErrorResponse{Message: err})
}

func InternalServer(w http.ResponseWriter, err string) {
	Error(w, http.StatusInternalServerError, err)
}

func NotFound(w http.ResponseWriter, err string) {
	Error(w, http.StatusNotFound, err)
}

func BadRequest(w http.ResponseWriter, err string) {
	Error(w, http.StatusBadRequest, err)
}

func Unauthorized(w http.ResponseWriter, err string) {
	Error(w, http.StatusUnauthorized, err)
}

func Forbidden(w http.ResponseWriter, err string) {
	Error(w, http.StatusForbidden, err)
}
