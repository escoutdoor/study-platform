package course

import (
	"net/http"

	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// delete godoc
//
//	@Summary		Delete course
//	@Description	Deletes a course owned by the authenticated teacher.
//	@Tags			courses
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int							true	"Course ID"
//	@Success		204	{string}	string						"No Content"
//	@Failure		400	{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401	{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	httpresponse.ErrorResponse	"Forbidden"
//	@Failure		404	{object}	httpresponse.ErrorResponse	"Course not found"
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/courses/{id} [delete]
func (h *handler) delete(w http.ResponseWriter, r *http.Request) error {
	courseID, err := validateCourseID(r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	userID, err := httpctx.GetID(ctx)
	if err != nil {
		return err
	}

	if err := h.service.Delete(ctx, courseID, userID); err != nil {
		return err
	}

	httpresponse.NoContent(w)
	return nil
}
