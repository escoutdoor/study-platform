package teacher

import (
	"net/http"

	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// get godoc
//
//	@Summary		Get teacher
//	@Description	Returns a teacher by user identifier.
//	@Tags			teachers
//	@Produce		json
//	@Param			id	path		int	true	"Teacher user ID"
//	@Success		200	{object}	getResponse
//	@Failure		400	{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		404	{object}	httpresponse.ErrorResponse	"Teacher not found"
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/teachers/{id} [get]
func (h *handler) get(w http.ResponseWriter, r *http.Request) error {
	userID, err := validateUserID(r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	teacher, err := h.service.Get(ctx, userID)
	if err != nil {
		return err
	}

	resp := getResponse{Teacher: teacherToResponse(teacher)}
	httpresponse.OK(w, resp)
	return nil
}

type getResponse struct {
	Teacher teacherResponse `json:"teacher"`
}
