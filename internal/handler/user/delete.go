package user

import (
	"net/http"

	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// delete godoc
//
//	@Summary		Delete current user
//	@Description	Deletes the authenticated user account.
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		204	{string}	string						"No Content"
//	@Failure		401	{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		404	{object}	httpresponse.ErrorResponse	"User not found"
//	@Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/users/me [delete]
func (h *handler) delete(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, err := httpctx.GetID(ctx)
	if err != nil {
		return err
	}

	if err := h.userService.Delete(ctx, userID); err != nil {
		return err
	}

	httpresponse.NoContent(w)
	return nil
}
