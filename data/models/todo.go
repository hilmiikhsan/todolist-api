package models

import "time"

type Todo struct {
	TodoID          int       `db:"id"`
	Title           string    `db:"title"`
	ActivityGroupID int       `db:"activity_group_id"`
	IsActive        bool      `db:"is_active"`
	Priority        string    `db:"priority"`
	UpdatedAt       time.Time `db:"updated_at"`
	CreatedAt       time.Time `db:"created_at"`
}
