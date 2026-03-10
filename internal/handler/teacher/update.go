package teacher

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
//	@Summary		Update current teacher
//	@Description	Updates the authenticated teacher profile.
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		updateRequest				true	"Update current teacher request"
//	@Success		200		{object}	updateResponse				"Teacher profile updated successfully"
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	httpresponse.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	httpresponse.ErrorResponse	"Teacher not found"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/teachers/me [put]
func (h *handler) update(w http.ResponseWriter, r *http.Request) error {
	req := new(updateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperror.ErrInvalidJSON
	}

	if err := h.cv.Validate(req); err != nil {
		return err
	}

	ctx := r.Context()
	userID, err := httpctx.GetID(ctx)
	if err != nil {
		return err
	}

	in := updateRequestToTeacher(req, userID)
	teacher, err := h.service.Update(ctx, in)
	if err != nil {
		return err
	}

	resp := updateResponse{Teacher: teacherToResponse(teacher)}
	httpresponse.OK(w, resp)
	return nil
}

type updateRequest struct {
	Department string `json:"department" validate:"required,min=2"`
}

func updateRequestToTeacher(req *updateRequest, userID int) entity.Teacher {
	return entity.Teacher{
		UserID:     userID,
		Department: req.Department,
	}
}

type updateResponse struct {
	Teacher teacherResponse `json:"teacher"`
}
