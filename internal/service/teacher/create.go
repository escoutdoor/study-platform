package teacher

import (
	"context"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Create(ctx context.Context, in entity.Teacher) error {
	if err := s.teacherRepo.Create(ctx, in); err != nil {
		return errwrap.Wrap("create teacher in repo", err)
	}

	return nil
}
