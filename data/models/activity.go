package models

import "time"

type Activity struct {
	ActivityID int       `db:"id"`
	Title      string    `db:"title"`
	Email      string    `db:"email"`
	UpdatedAt  time.Time `db:"updated_at"`
	CreatedAt  time.Time `db:"created_at"`
}
