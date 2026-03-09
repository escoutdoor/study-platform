package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/internal/entity"
	"github.com/escoutdoor/study-platform/internal/util/token"
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

func (r *Repository) GetByID(ctx context.Context, userID int) (entity.User, error) {
	sqlStatement := `
		SELECT 
            id,
            first_name,
            last_name,
            email,
            password,
            created_at,
            updated_at
		FROM users
		WHERE id = $1
    `
	q := database.Query{
		Name: "user_repository.GetByID",
		Sql:  sqlStatement,
	}

	var user entity.User
	if err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, apperror.UserNotFoundID(userID)
		}

		return entity.User{}, scanRowError(err)
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	sqlStatement := `
		SELECT 
            id,
            first_name,
            last_name,
            email,
            password,
            created_at,
            updated_at
		FROM users
		WHERE email = $1
    `
	q := database.Query{
		Name: "user_repository.GetByEmail",
		Sql:  sqlStatement,
	}

	var user entity.User
	if err := r.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, apperror.UserNotFoundEmail(email)
		}

		return entity.User{}, scanRowError(err)
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, in entity.User) (int, error) {
	sql := `
        INSERT INTO users(first_name,last_name,email,password)
        VALUES($1,$2,$3,$4)
        RETURNING id
    `

	q := database.Query{
		Name: "user_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.FirstName,
		in.LastName,
		in.Email,
		in.Password,
	}

	var userID int
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(&userID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				return 0, apperror.UserEmailAlreadyExists(in.Email)
			}
		}

		return 0, scanRowError(err)
	}

	return userID, nil
}

func (r *Repository) Update(ctx context.Context, in entity.User) (entity.User, error) {
	sql := `
        UPDATE users SET
            first_name=$1,
            last_name=$2,
            email=$3,
            password=$4,
            updated_at=now()
        WHERE id=$5
        RETURNING id,first_name,last_name,email,created_at,updated_at
    `

	q := database.Query{
		Name: "user_repository.Update",
		Sql:  sql,
	}

	args := []any{
		in.FirstName,
		in.LastName,
		in.Email,
		in.Password,
		in.ID,
	}

	var user entity.User
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				return entity.User{}, apperror.UserEmailAlreadyExists(in.Email)
			}
		}
		return entity.User{}, scanRowError(err)
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, userID int) error {
	sql := `
		DELETE FROM users WHERE id=$1
	`

	q := database.Query{
		Name: "user_repository.Delete",
		Sql:  sql,
	}

	cmd, err := r.db.ExecContext(ctx, q, userID)
	if err != nil {
		return executeSQLError(err)
	}

	if v, _ := cmd.RowsAffected(); v == 0 {
		return fmt.Errorf("delete rows affected: 0")
	}

	return nil
}

func (r *Repository) GetRoles(ctx context.Context, userID int) ([]token.Role, error) {
	sqlStatement := `
		SELECT 
			(EXISTS(SELECT 1 FROM students WHERE user_id = $1)) AS is_student,
			(EXISTS(SELECT 1 FROM teachers WHERE user_id = $1)) AS is_teacher
	`
	q := database.Query{
		Name: "user_repository.GetRoles",
		Sql:  sqlStatement,
	}

	var isStudent, isTeacher bool
	if err := r.db.QueryRowContext(ctx, q, userID).Scan(&isStudent, &isTeacher); err != nil {
		return nil, scanRowError(err)
	}

	var roles []token.Role
	if isStudent {
		roles = append(roles, token.RoleStudent)
	}
	if isTeacher {
		roles = append(roles, token.RoleTeacher)
	}

	return roles, nil
}
