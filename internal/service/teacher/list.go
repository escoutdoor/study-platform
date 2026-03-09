package teacher

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) List(ctx context.Context) ([]entity.Teacher, error) {
	teachers, err := s.teacherRepo.List(ctx)
	if err != nil {
		return nil, errwrap.Wrap("get list of teachers from repo", err)
	}

	return teachers, nil
}
