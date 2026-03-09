package student

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Update(ctx context.Context, in entity.Student) error {
	if _, err := s.studentRepo.GetByUserID(ctx, in.UserID); err != nil {
		return errwrap.Wrap("get student from repo by id", err)
	}

	if err := s.studentRepo.Update(ctx, in); err != nil {
		return errwrap.Wrap("update student in repo", err)
	}

	return nil
}
