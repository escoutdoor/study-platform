package student

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) List(ctx context.Context) ([]entity.Student, error) {
	students, err := s.studentRepo.List(ctx)
	if err != nil {
		return nil, errwrap.Wrap("get list of students from repo", err)
	}

	return students, nil
}
