package teacher

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
)

type Service struct {
	teacherRepo teacherRepository
}

func New(teacherRepo teacherRepository) *Service {
	return &Service{
		teacherRepo: teacherRepo,
	}
}

type teacherRepository interface {
	List(ctx context.Context) ([]entity.Teacher, error)
	GetByUserID(ctx context.Context, userID int) (entity.Teacher, error)

	Create(ctx context.Context, in entity.Teacher) error
	Update(ctx context.Context, in entity.Teacher) error
}
