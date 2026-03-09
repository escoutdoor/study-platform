package student

import (
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
	"net/http"
)

// get godoc
//
//	@Summary		Get student
//	@Description	Returns a student by user identifier.
//	@Tags			students
//	@Produce		json
//	@Param			id	path		int	true	"Student user ID"
//	@Success		200	{object}	getResponse
//	@Failure		400	{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		404	{object}	httpresponse.ErrorResponse	"Student not found"
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/students/{id} [get]
func (h *handler) get(w http.ResponseWriter, r *http.Request) error {
	userID, err := validateUserID(r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	student, err := h.studentService.Get(ctx, userID)
	if err != nil {
		return err
	}

	resp := getResponse{Student: studentToResponse(student)}
	httpresponse.OK(w, resp)
	return nil
}

type getResponse struct {
	Student studentResponse `json:"student"`
}
