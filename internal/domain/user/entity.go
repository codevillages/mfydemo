package user

import "time"

// User represents the domain entity.
type User struct {
	ID         int64
	Username   string
	Password   string // hashed password
	Email      string
	Status     int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
