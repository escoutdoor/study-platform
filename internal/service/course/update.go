package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Update(ctx context.Context, in entity.Course) (entity.Course, error) {
	course, err := s.courseRepo.Get(ctx, in.ID)
	if err != nil {
		return entity.Course{}, errwrap.Wrap("get course from repo by id", err)
	}
	if course.TeacherID != in.TeacherID {
		return entity.Course{}, apperror.CourseAccessDenied
	}

	updatedCourse, err := s.courseRepo.Update(ctx, in)
	if err != nil {
		return entity.Course{}, errwrap.Wrap("update course in repo", err)
	}

	return updatedCourse, nil
}
