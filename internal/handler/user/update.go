package user

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// updateMe godoc
//
//	@Summary		Update current user
//	@Description	Updates the authenticated user's profile data and returns the updated user.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		updateRequest	true	"Update current user request"
//	@Success		200		{object}	updateResponse
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		404		{object}	httpresponse.ErrorResponse	"User not found"
//	@Failure		409		{object}	httpresponse.ErrorResponse	"User with this email already exists"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/users/me [put]
func (h *handler) updateMe(w http.ResponseWriter, r *http.Request) error {
	req := new(updateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperror.ErrInvalidJSON
	}

	if err := h.cv.Validate(req); err != nil {
		return err
	}

	ctx := r.Context()
	userID, err := httpctx.GetID(ctx)
	if err != nil {
		return err
	}

	in := updateRequestToUser(req, userID)
	user, err := h.userService.Update(ctx, in)
	if err != nil {
		return err
	}

	resp := updateResponse{User: userToResponse(user)}
	httpresponse.OK(w, resp)
	return nil
}

type updateRequest struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=20"`
	LastName  string `json:"lastName" validate:"required,min=1,max=20"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func updateRequestToUser(req *updateRequest, id int) entity.User {
	return entity.User{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,

		Email:    req.Email,
		Password: req.Password,
	}
}

type updateResponse struct {
	User userResponse `json:"user"`
}
