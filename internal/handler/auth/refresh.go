package auth

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// refreshToken godoc
//
//	@Summary		Refresh tokens
//	@Description	Refreshes access and refresh tokens using a refresh token.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		refreshTokenRequest	true	"Refresh token request"
//	@Success		200		{object}	refreshTokenResponse
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Invalid or expired refresh token"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/auth/refresh [post]
func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) error {
	req := new(refreshTokenRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperror.ErrInvalidJSON
	}

	if err := h.cv.Validate(req); err != nil {
		return err
	}

	ctx := r.Context()
	tokens, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return err
	}

	resp := refreshTokenResponse{Tokens: tokensToResponse(tokens)}
	httpresponse.OK(w, resp)
	return nil
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type refreshTokenResponse struct {
	Tokens authResponse `json:"tokens"`
}
