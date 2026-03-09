package course

import (
	"context"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Create(ctx context.Context, in entity.Course) (int, error) {
	courseID, err := s.courseRepo.Create(ctx, in)
	if err != nil {
		return 0, errwrap.Wrap("create course in repo", err)
	}

	return courseID, nil
}
