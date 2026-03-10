package teacher

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Update(ctx context.Context, in entity.Teacher) (entity.Teacher, error) {
	if _, err := s.teacherRepo.GetByUserID(ctx, in.UserID); err != nil {
		return entity.Teacher{}, errwrap.Wrap("get teacher from repo by id", err)
	}

	if err := s.teacherRepo.Update(ctx, in); err != nil {
		return entity.Teacher{}, errwrap.Wrap("update teacher in repo", err)
	}

	updatedTeacher, err := s.teacherRepo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return entity.Teacher{}, errwrap.Wrap("get teacher from repo by id", err)
	}

	return updatedTeacher, nil
}
