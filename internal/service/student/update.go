package student

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Update(ctx context.Context, in entity.Student) (entity.Student, error) {
	if _, err := s.studentRepo.GetByUserID(ctx, in.UserID); err != nil {
		return entity.Student{}, errwrap.Wrap("get student from repo by id", err)
	}

	if err := s.studentRepo.Update(ctx, in); err != nil {
		return entity.Student{}, errwrap.Wrap("update student in repo", err)
	}

	updatedStudent, err := s.studentRepo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return entity.Student{}, errwrap.Wrap("get student from repo by id", err)
	}

	return updatedStudent, nil
}
