package teacher

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Get(ctx context.Context, userID int) (entity.Teacher, error) {
	teacher, err := s.teacherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return entity.Teacher{}, errwrap.Wrap("get teacher from repo", err)
	}

	return teacher, nil
}
