package user

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
)

type Service struct {
	userRepo userRepository
}

func New(userRepo userRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

type userRepository interface {
	GetByID(ctx context.Context, userID int) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)

	Update(ctx context.Context, in entity.User) (entity.User, error)
	Delete(ctx context.Context, userID int) error
}
