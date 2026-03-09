package student

import (
	"net/http"

	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// enroll godoc
//
//	@Summary		Enroll in course
//	@Description	Enrolls the authenticated student in a course.
//	@Tags			students
//	@Produce		json
//	@Security		BearerAuth
//	@Param			courseId	path		int	true	"Course ID"
//	@Success		201		{string}	string	"Created"
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	httpresponse.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	httpresponse.ErrorResponse	"Course not found"
//	@Failure		409		{object}	httpresponse.ErrorResponse	"Student is already enrolled in this course"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/students/me/courses/{courseId} [post]
func (h *handler) enroll(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, err := httpctx.GetID(ctx)
	if err != nil {
		return err
	}

	courseID, err := validateCourseID(r)
	if err != nil {
		return err
	}

	if err := h.enrollmentService.Enroll(ctx, userID, courseID); err != nil {
		return err
	}

	httpresponse.Created(w, nil)
	return nil
}
