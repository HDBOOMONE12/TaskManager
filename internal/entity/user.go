package entity

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
