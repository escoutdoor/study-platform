package course

import (
	"net/http"

	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// list godoc
//
//	@Summary		List courses
//	@Description	Returns the list of all courses.
//	@Tags			courses
//	@Produce		json
//	@Success		200	{object}	listResponse
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/courses [get]
func (h *handler) list(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	list, err := h.service.List(ctx)
	if err != nil {
		return err
	}

	resp := listResponse{Courses: courseListToResponse(list)}
	httpresponse.OK(w, resp)
	return nil
}

type listResponse struct {
	Courses []courseResponse `json:"courses"`
}
