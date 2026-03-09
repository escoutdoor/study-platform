package student

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

func (r *Repository) List(ctx context.Context) ([]entity.Student, error) {
	sql := `
		SELECT 
			s.user_id, 
			u.first_name, 
			u.last_name, 
			u.email, 
			s.created_at,
            s.updated_at
		FROM students s
		JOIN users u ON s.user_id = u.id
	`
	q := database.Query{
		Name: "student_repository.List",
		Sql:  sql,
	}

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var students []entity.Student
	for rows.Next() {
		var s entity.Student
		if err := rows.Scan(
			&s.UserID,
			&s.FirstName,
			&s.LastName,
			&s.Email,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, scanRowsError(err)
		}

		students = append(students, s)
	}

	return students, nil
}

func (r *Repository) GetByUserID(ctx context.Context, studentID int) (entity.Student, error) {
	sqlStatement := `
		SELECT 
			s.user_id, 
			u.first_name, 
			u.last_name, 
			u.email, 
			s.created_at,
            s.updated_at
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.user_id = $1
    `
	q := database.Query{
		Name: "student_repository.GetByUserID",
		Sql:  sqlStatement,
	}

	var student entity.Student
	if err := r.db.QueryRowContext(ctx, q, studentID).Scan(
		&student.UserID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.CreatedAt,
		&student.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Student{}, apperror.StudentNotFoundID(studentID)
		}

		return entity.Student{}, scanRowError(err)
	}

	return student, nil
}

func (r *Repository) Create(ctx context.Context, in entity.Student) error {
	sql := `
        INSERT INTO students(user_id)
        VALUES($1)
    `

	q := database.Query{
		Name: "student_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.UserID,
	}

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, in entity.Student) error {
	sql := `
        UPDATE students SET
            updated_at=now()
        WHERE user_id=$1
    `

	q := database.Query{
		Name: "student_repository.Update",
		Sql:  sql,
	}

	if _, err := r.db.ExecContext(ctx, q, in.UserID); err != nil {
		return executeSQLError(err)
	}

	return nil
}
