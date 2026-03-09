package entity

import "time"

type Course struct {
	ID        int
	TeacherID int

	Title       string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}
