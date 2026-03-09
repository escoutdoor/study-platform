package student

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
)

type Service struct {
	studentRepo studentRepository
}

func New(studentRepo studentRepository) *Service {
	return &Service{
		studentRepo: studentRepo,
	}
}

type studentRepository interface {
	List(ctx context.Context) ([]entity.Student, error)
	GetByUserID(ctx context.Context, studentID int) (entity.Student, error)

	Update(ctx context.Context, in entity.Student) error
}
