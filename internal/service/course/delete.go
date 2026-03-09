package course

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Delete(ctx context.Context, courseID, teacherID int) error {
	course, err := s.courseRepo.Get(ctx, courseID)
	if err != nil {
		return errwrap.Wrap("get course from repo", err)
	}
	if course.TeacherID != teacherID {
		return apperror.CourseAccessDenied
	}

	if err := s.courseRepo.Delete(ctx, courseID); err != nil {
		return errwrap.Wrap("delete course in repo", err)
	}

	return nil
}
