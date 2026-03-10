package teacher

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/httpctx"
	"github.com/escoutdoor/study-platform/pkg/httpresponse"
)

// create godoc
//
//	@Summary		Create teacher profile
//	@Description	Creates a teacher profile for the authenticated user.
//	@Tags			teachers
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createRequest				true	"Create teacher request"
//	@Success		201		{object}	map[string]string			"Teacher profile created successfully"
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		409		{object}	httpresponse.ErrorResponse	"Teacher already exists"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/teachers [post]
func (h *handler) create(w http.ResponseWriter, r *http.Request) error {
	req := new(createRequest)
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

	in := createRequestToTeacher(req, userID)
	if err := h.service.Create(ctx, in); err != nil {
		return err
	}

	httpresponse.Created(w, map[string]string{"message": "teacher profile created successfully"})
	return nil
}

type createRequest struct {
	Department string `json:"department" validate:"required,min=3"`
}

func createRequestToTeacher(req *createRequest, userID int) entity.Teacher {
	return entity.Teacher{
		UserID:     userID,
		Department: req.Department,
	}
}
