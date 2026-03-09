package entity

import "time"

type Teacher struct {
	UserID int

	FirstName string
	LastName  string

	Department string

	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}
