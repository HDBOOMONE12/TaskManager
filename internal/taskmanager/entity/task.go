package entity

import "time"

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Status      string
	Priority    int64
	DueAt       *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
