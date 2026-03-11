package auth

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// login godoc
//
//	@Summary		Login user
//	@Description	Authenticates a user with email and password and returns access and refresh tokens.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Incorrect credentials"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (h *handler) login(w http.ResponseWriter, r *http.Request) error {
	req := new(loginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperror.ErrInvalidJSON
	}

	if err := h.cv.Validate(req); err != nil {
		return err
	}

	ctx := r.Context()
	in := loginRequestToUser(req)

	tokens, err := h.authService.Login(ctx, in)
	if err != nil {
		return err
	}

	resp := loginResponse{Tokens: tokensToResponse(tokens)}
	httpresponse.OK(w, resp)
	return nil
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}

func loginRequestToUser(req *loginRequest) entity.User {
	return entity.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

type loginResponse struct {
	Tokens authResponse `json:"tokens"`
}
