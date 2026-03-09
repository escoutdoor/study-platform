package user

import (
	"context"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/errhandler"
	"github.com/escoutdoor/study-platform/pkg/validator"
)

type handler struct {
	userService userService
	cv          *validator.CustomValidator
}

func RegisterHandlers(
	mux *http.ServeMux,
	userService userService,
	cv *validator.CustomValidator,
	authMiddleware func(errhandler.HandlerFunc) errhandler.HandlerFunc,
) {
	h := &handler{
		userService: userService,
		cv:          cv,
	}

	routes := map[string]errhandler.HandlerFunc{
		"PUT /users/me": authMiddleware(h.updateMe),

		"DELETE /users/me": authMiddleware(h.delete),
	}
	for p, h := range routes {
		mux.Handle(p, errhandler.ErrorHandler(h))
	}
}

type userService interface {
	Update(ctx context.Context, in entity.User) (entity.User, error)
	Delete(ctx context.Context, userID int) error
}

type userResponse struct {
	ID int `json:"id"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email string `json:"email"`
}

func userToResponse(user entity.User) userResponse {
	return userResponse{
		ID: user.ID,

		FirstName: user.FirstName,
		LastName:  user.LastName,

		Email: user.Email,
	}
}
