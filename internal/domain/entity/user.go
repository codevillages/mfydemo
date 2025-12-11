package entity

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	Status    int32
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
