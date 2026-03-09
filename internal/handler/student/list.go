package student

import (
	"net/http"

	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// list godoc
//
//	@Summary		List students
//	@Description	Returns the list of all students.
//	@Tags			students
//	@Produce		json
//	@Success		200	{object}	listResponse
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/students [get]
func (h *handler) list(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	list, err := h.studentService.List(ctx)
	if err != nil {
		return err
	}

	resp := listResponse{Students: studentListToResponse(list)}
	httpresponse.OK(w, resp)
	return nil
}

type listResponse struct {
	Students []studentResponse `json:"students"`
}
