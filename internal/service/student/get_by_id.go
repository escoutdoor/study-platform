package student

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Get(ctx context.Context, studentID int) (entity.Student, error) {
	student, err := s.studentRepo.GetByUserID(ctx, studentID)
	if err != nil {
		return entity.Student{}, errwrap.Wrap("get student from repo", err)
	}

	return student, nil
}
