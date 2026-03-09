package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) List(ctx context.Context) ([]entity.Course, error) {
	courses, err := s.courseRepo.List(ctx)
	if err != nil {
		return nil, errwrap.Wrap("get list of courses from repo", err)
	}

	return courses, nil
}
