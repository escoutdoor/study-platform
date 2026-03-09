package teacher

import (
	"net/http"

	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// list godoc
//
//	@Summary		List teachers
//	@Description	Returns the list of all teachers.
//	@Tags			teachers
//	@Produce		json
//	@Success		200	{object}	listResponse
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/teachers [get]
func (h *handler) list(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	list, err := h.service.List(ctx)
	if err != nil {
		return err
	}

	resp := listResponse{Teachers: teacherListToResponse(list)}
	httpresponse.OK(w, resp)
	return nil
}

type listResponse struct {
	Teachers []teacherResponse `json:"teachers"`
}
