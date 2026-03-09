package teacher

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

func (r *Repository) List(ctx context.Context) ([]entity.Teacher, error) {
	sql := `
		SELECT 
			t.user_id, 
            t.department,
			u.first_name, 
			u.last_name, 
			u.email, 
			t.created_at,
            t.updated_at
		FROM teachers t
		JOIN users u ON t.user_id = u.id
	`
	q := database.Query{
		Name: "teacher_repository.List",
		Sql:  sql,
	}

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var teachers []entity.Teacher
	for rows.Next() {
		var t entity.Teacher
		if err := rows.Scan(
			&t.UserID,
			&t.Department,
			&t.FirstName,
			&t.LastName,
			&t.Email,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, scanRowsError(err)
		}

		teachers = append(teachers, t)
	}

	return teachers, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int) (entity.Teacher, error) {
	sqlStatement := `
		SELECT 
			t.user_id, 
            t.department,
			u.first_name, 
			u.last_name, 
			u.email, 
			t.created_at,
            t.updated_at
		FROM teachers t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1
    `
	q := database.Query{
		Name: "teacher_repository.GetByUserID",
		Sql:  sqlStatement,
	}

	var teacher entity.Teacher
	if err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&teacher.UserID,
		&teacher.Department,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Teacher{}, apperror.TeacherNotFoundID(userID)
		}

		return entity.Teacher{}, scanRowError(err)
	}

	return teacher, nil
}

func (r *Repository) Create(ctx context.Context, in entity.Teacher) error {
	sql := `
        INSERT INTO teachers(user_id, department)
        VALUES($1,$2)
    `

	q := database.Query{
		Name: "teacher_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.UserID,
		in.Department,
	}

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return apperror.TeacherAlreadyExists
			}
		}

		return executeSQLError(err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, in entity.Teacher) error {
	sql := `
        UPDATE teachers SET
            department=$1,
            updated_at=now()
        WHERE user_id=$2
    `

	q := database.Query{
		Name: "teacher_repository.Update",
		Sql:  sql,
	}

	args := []any{
		in.Department,
		in.UserID,
	}

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		return executeSQLError(err)
	}

	return nil
}
