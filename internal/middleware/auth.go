package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/internal/util/token"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

const (
	authorizationHeader       = "Authorization"
	authorizationHeaderPrefix = "Bearer "
)

type tokenProvider interface {
	ValidateAccessToken(accessToken string) (token.AccessTokenClaims, error)
}

func Auth(tp tokenProvider) func(errhandler.HandlerFunc) errhandler.HandlerFunc {
	return func(next errhandler.HandlerFunc) errhandler.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			authHeader := r.Header.Get(authorizationHeader)
			if authHeader == "" {
				httpresponse.Unauthorized(w, "authorization header not provided")
				return nil
			}

			if !strings.HasPrefix(authHeader, authorizationHeaderPrefix) {
				httpresponse.Unauthorized(w, "invalid authorization header format")
				return nil
			}

			token := authHeader[len(authorizationHeaderPrefix):]

			claims, err := tp.ValidateAccessToken(token)
			if err != nil {
				httpresponse.Unauthorized(w, "invalid or expired token")
				return nil
			}

			ctx := context.WithValue(r.Context(), httpctx.UserIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, httpctx.RolesContextKey, claims.Roles)

			r = r.WithContext(ctx)

			return next(w, r)
		}
	}
}

func RequireRole(allowedRoles ...token.Role) func(errhandler.HandlerFunc) errhandler.HandlerFunc {
	return func(next errhandler.HandlerFunc) errhandler.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			roles, err := httpctx.GetRoles(ctx)
			if err != nil {
				httpresponse.Forbidden(w, "unauthorized or missing roles")
				return nil
			}

			hasAccess := false
			for _, allowedRole := range allowedRoles {
				if slices.Contains(roles, allowedRole) {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				httpresponse.Forbidden(w, "insufficient permissions")
				return nil
			}

			return next(w, r)
		}
	}
}
