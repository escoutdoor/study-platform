package entity

import "time"

type Student struct {
	UserID int

	FirstName string
	LastName  string

	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}
