package entity

import "time"

type User struct {
	ID int

	FirstName string
	LastName  string

	Email    string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
}
