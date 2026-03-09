package auth

import (
	"context"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/pkg/validator"
)

type handler struct {
	authService authService
	cv          *validator.CustomValidator
}

func RegisterHandlers(
	mux *http.ServeMux,
	authService authService,
	cv *validator.CustomValidator,
) {
	h := &handler{
		authService: authService,
		cv:          cv,
	}

	routes := map[string]errhandler.HandlerFunc{
		"POST /auth/register": h.register,
		"POST /auth/login":    h.login,

		"POST /auth/refresh": h.refreshToken,
	}
	for p, h := range routes {
		mux.Handle(p, errhandler.ErrorHandler(h))
	}
}

type authService interface {
	Register(ctx context.Context, in entity.User) (entity.Tokens, error)
	Login(ctx context.Context, in entity.User) (entity.Tokens, error)

	RefreshToken(ctx context.Context, refreshToken string) (entity.Tokens, error)
}

type authResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func tokensToResponse(tokens entity.Tokens) authResponse {
	return authResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
