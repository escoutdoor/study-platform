package enrollment

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/database"
	"github.com/lib/pq"
)

type Repository struct {
	db database.DB
}

func New(db database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Get(ctx context.Context, studentID, courseID int) (entity.Enrollment, error) {
	sqlStatement := `
        SELECT 
            student_id,
            course_id 
        FROM enrollments 
        WHERE student_id=$1 AND course_id=$2
    `
	q := database.Query{
		Name: "enrollment_repository.Get",
		Sql:  sqlStatement,
	}

	var enrollment entity.Enrollment
	if err := r.db.QueryRowContext(ctx, q, studentID, courseID).Scan(
		&enrollment.StudentID,
		&enrollment.CourseID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Enrollment{}, apperror.StudentNotEnrolled
		}

		return entity.Enrollment{}, scanRowError(err)
	}

	return enrollment, nil
}

func (r *Repository) Create(ctx context.Context, studentID, courseID int) error {
	sql := `
        INSERT INTO enrollments(student_id,course_id)
        VALUES($1,$2)
    `

	q := database.Query{
		Name: "enrollment_repository.Create",
		Sql:  sql,
	}

	args := []any{
		studentID,
		courseID,
	}
	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return apperror.StudentAlreadyEnrolled
			}
		}

		return executeSQLError(err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, studentID, courseID int) error {
	sql := `
		DELETE 
            FROM enrollments 
        WHERE 
            student_id=$1 AND course_id=$2
	`

	q := database.Query{
		Name: "enrollment_repository.Delete",
		Sql:  sql,
	}

	args := []any{
		studentID,
		courseID,
	}

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		return executeSQLError(err)
	}

	return nil
}
