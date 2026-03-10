package enrollment

import (
	"context"

	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

func (s *Service) Enroll(ctx context.Context, studentID, courseID int) error {
	if _, err := s.courseRepo.Get(ctx, courseID); err != nil {
		return errwrap.Wrap("get course from course repo", err)
	}

	if err := s.enrollmentRepo.Create(ctx, studentID, courseID); err != nil {
		return errwrap.Wrap("create enrollment in enrollment repo", err)
	}

	return nil
}
