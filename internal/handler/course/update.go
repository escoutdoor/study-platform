package course

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// update godoc
//
//	@Summary		Update course
//	@Description	Updates a course owned by the authenticated teacher.
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int				true	"Course ID"
//	@Param			request	body		updateRequest	true	"Update course request"
//	@Success		200		{object}	updateResponse
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	httpresponse.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	httpresponse.ErrorResponse	"Course not found"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/courses/{id} [put]
func (h *handler) update(w http.ResponseWriter, r *http.Request) error {
	req := new(updateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperror.ErrInvalidJSON
	}

	if err := h.cv.Validate(req); err != nil {
		return err
	}

	courseID, err := validateCourseID(r)
	if err != nil {
		return err
	}

	ctx := r.Context()
	userID, err := httpctx.GetID(ctx)
	if err != nil {
		return err
	}

	in := updateRequestToCourse(req, courseID, userID)
	course, err := h.service.Update(ctx, in)
	if err != nil {
		return err
	}

	resp := updateResponse{Course: courseToResponse(course)}
	httpresponse.OK(w, resp)
	return nil
}

type updateRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3"`
}

func updateRequestToCourse(req *updateRequest, courseID int, userID int) entity.Course {
	return entity.Course{
		ID:        courseID,
		TeacherID: userID,

		Title:       req.Title,
		Description: req.Description,
	}
}

type updateResponse struct {
	Course courseResponse `json:"course"`
}
