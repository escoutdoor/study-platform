package student

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
//	@Summary		Update current student
//	@Description	Updates the authenticated student. The request schema is temporary and may change!!!
//	@Tags			students
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		updateRequest				true	"Temporary update current student request"
//	@Success		200		{object}	updateResponse				"Student profile updated successfully"
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	httpresponse.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	httpresponse.ErrorResponse	"Student not found"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/students/me [put]
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

	in := updateRequestToStudent(req, userID)
	student, err := h.studentService.Update(ctx, in)
	if err != nil {
		return err
	}

	resp := updateResponse{Student: studentToResponse(student)}
	httpresponse.OK(w, resp)
	return nil
}

type updateRequest struct {
	// TODO: there should be some fields to update
}

func updateRequestToStudent(req *updateRequest, userID int) entity.Student {
	return entity.Student{
		UserID: userID,
	}
}

type updateResponse struct {
	Student studentResponse `json:"student"`
}
