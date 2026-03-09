package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Get(ctx context.Context, courseID int) (entity.Course, error) {
	course, err := s.courseRepo.Get(ctx, courseID)
	if err != nil {
		return entity.Course{}, errwrap.Wrap("get course from repo", err)
	}

	return course, nil
}
