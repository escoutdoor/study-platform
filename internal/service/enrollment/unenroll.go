package enrollment

import (
	"context"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Unenroll(ctx context.Context, studentID, courseID int) error {
	if _, err := s.courseRepo.Get(ctx, courseID); err != nil {
		return errwrap.Wrap("get course from course repo", err)
	}

	if _, err := s.enrollmentRepo.Get(ctx, studentID, courseID); err != nil {
		return errwrap.Wrap("get enrollment from enrollment repo", err)
	}

	if err := s.enrollmentRepo.Delete(ctx, studentID, courseID); err != nil {
		return errwrap.Wrap("delete enrollment in enrollment repo", err)
	}

	return nil
}
