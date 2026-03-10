package course

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
//	@Summary		Create course
//	@Description	Creates a new course for the authenticated teacher.
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createRequest	true	"Create course request"
//	@Success		201		{object}	createResponse
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		401		{object}	httpresponse.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	httpresponse.ErrorResponse	"Forbidden"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/courses [post]
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
	in := createRequestToCourse(req, userID)

	courseID, err := h.service.Create(ctx, in)
	if err != nil {
		return err
	}

	resp := createResponse{CourseID: courseID}
	httpresponse.Created(w, resp)
	return nil
}

type createRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3"`
}

func createRequestToCourse(req *createRequest, userID int) entity.Course {
	return entity.Course{
		TeacherID: userID,

		Title:       req.Title,
		Description: req.Description,
	}
}

type createResponse struct {
	CourseID int `json:"courseId"`
}
