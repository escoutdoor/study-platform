package course

import (
	"net/http"

	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// get godoc
//
//	@Summary		Get course
//	@Description	Returns a course by its identifier.
//	@Tags			courses
//	@Produce		json
//	@Param			id	path		int	true	"Course ID"
//	@Success		200	{object}	getResponse
//	@Failure		400	{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		404	{object}	httpresponse.ErrorResponse	"Course not found"
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/courses/{id} [get]
func (h *handler) get(w http.ResponseWriter, r *http.Request) error {
	courseID, err := validateCourseID(r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	course, err := h.service.Get(ctx, courseID)
	if err != nil {
		return err
	}

	resp := getResponse{Course: courseToResponse(course)}
	httpresponse.OK(w, resp)
	return nil
}

type getResponse struct {
	Course courseResponse `json:"course"`
}
