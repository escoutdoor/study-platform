package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
)

type Service struct {
	courseRepo courseRepository
}

func New(courseRepo courseRepository) *Service {
	return &Service{
		courseRepo: courseRepo,
	}
}

type courseRepository interface {
	List(ctx context.Context) ([]entity.Course, error)
	Get(ctx context.Context, courseID int) (entity.Course, error)

	Create(ctx context.Context, in entity.Course) (int, error)
	Update(ctx context.Context, in entity.Course) (entity.Course, error)
	Delete(ctx context.Context, courseID int) error
}
