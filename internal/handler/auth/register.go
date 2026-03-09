package auth

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// register godoc
//
//	@Summary		Register user
//	@Description	Registers a new user account and returns access and refresh tokens.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		registerRequest	true	"Register request"
//	@Success		201		{object}	registerResponse
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		409		{object}	httpresponse.ErrorResponse	"User with this email already exists"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/auth/register [post]
func (h *handler) register(w http.ResponseWriter, r *http.Request) error {
	req := new(registerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperror.ErrInvalidJSON
	}

	if err := h.cv.Validate(req); err != nil {
		return err
	}

	ctx := r.Context()
	in := registerRequestToUser(req)

	tokens, err := h.authService.Register(ctx, in)
	if err != nil {
		return err
	}

	resp := registerResponse{Tokens: tokensToResponse(tokens)}
	httpresponse.Created(w, resp)
	return nil
}

type registerRequest struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=20"`
	LastName  string `json:"lastName" validate:"required,min=1,max=20"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func registerRequestToUser(req *registerRequest) entity.User {
	return entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,

		Email:    req.Email,
		Password: req.Password,
	}
}

type registerResponse struct {
	Tokens authResponse `json:"tokens"`
}
