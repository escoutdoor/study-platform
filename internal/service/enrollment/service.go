package enrollment

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
)

type Service struct {
	enrollmentRepo enrollmentRepository
	courseRepo     courseRepository
}

func New(enrollmentRepo enrollmentRepository, courseRepo courseRepository) *Service {
	return &Service{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
	}
}

type enrollmentRepository interface {
	Create(ctx context.Context, studentID, courseID int) error
	Delete(ctx context.Context, studentID, courseID int) error
	Get(ctx context.Context, studentID, courseID int) (entity.Enrollment, error)
}

type courseRepository interface {
	Get(ctx context.Context, courseID int) (entity.Course, error)
}
