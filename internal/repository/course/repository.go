package course

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/pkg/database"
)

type Repository struct {
	db database.DB
}

func New(db database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context) ([]entity.Course, error) {
	sql := `
        SELECT 
            id,
            teacher_id,
            title,
            description,
            created_at,
            updated_at
        FROM courses 
    `
	q := database.Query{
		Name: "course_repository.List",
		Sql:  sql,
	}

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var courses []entity.Course
	for rows.Next() {
		var c entity.Course
		if err := rows.Scan(
			&c.ID,
			&c.TeacherID,
			&c.Title,
			&c.Description,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, scanRowsError(err)
		}

		courses = append(courses, c)
	}

	return courses, nil
}

func (r *Repository) Get(ctx context.Context, courseID int) (entity.Course, error) {
	sqlStatement := `
        SELECT 
            id,
            teacher_id,
            title,
            description,
            created_at,
            updated_at
        FROM courses 
        WHERE id=$1
    `
	q := database.Query{
		Name: "course_repository.Get",
		Sql:  sqlStatement,
	}

	var course entity.Course
	if err := r.db.QueryRowContext(ctx, q, courseID).Scan(
		&course.ID,
		&course.TeacherID,
		&course.Title,
		&course.Description,
		&course.CreatedAt,
		&course.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Course{}, apperror.CourseNotFoundID(courseID)
		}

		return entity.Course{}, scanRowError(err)
	}

	return course, nil
}

func (r *Repository) Create(ctx context.Context, in entity.Course) (int, error) {
	sql := `
        INSERT INTO courses(teacher_id,title,description)
        VALUES($1,$2,$3)
        RETURNING id
    `

	q := database.Query{
		Name: "course_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.TeacherID,
		in.Title,
		in.Description,
	}

	var courseID int
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(
		&courseID,
	); err != nil {
		return 0, scanRowError(err)
	}

	return courseID, nil
}

func (r *Repository) Update(ctx context.Context, in entity.Course) (entity.Course, error) {
	sql := `
        UPDATE courses SET
            teacher_id=$1,
            title=$2,
            description=$3,
            updated_at=now()
        WHERE id=$4
        RETURNING id,teacher_id,title,description,created_at,updated_at
    `

	q := database.Query{
		Name: "course_repository.Update",
		Sql:  sql,
	}

	args := []any{
		in.TeacherID,
		in.Title,
		in.Description,
		in.ID,
	}

	var course entity.Course
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(
		&course.ID,
		&course.TeacherID,
		&course.Title,
		&course.Description,
		&course.CreatedAt,
		&course.UpdatedAt,
	); err != nil {
		return entity.Course{}, scanRowError(err)
	}

	return course, nil
}

func (r *Repository) Delete(ctx context.Context, courseID int) error {
	sql := `
		DELETE FROM courses WHERE id=$1
	`

	q := database.Query{
		Name: "course_repository.Delete",
		Sql:  sql,
	}

	if _, err := r.db.ExecContext(ctx, q, courseID); err != nil {
		return executeSQLError(err)
	}

	return nil
}
