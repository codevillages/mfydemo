package entity

import "github.com/gogf/gf/v2/os/gtime"

// User 表映射。
type User struct {
	Id           uint64      `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	Nickname     string      `json:"nickname"`
	Avatar       string      `json:"avatar"`
	PasswordHash string      `json:"-"`
	Status       int         `json:"status"`
	CreatedAt    *gtime.Time `json:"created_at"`
	UpdatedAt    *gtime.Time `json:"updated_at"`
	DeletedAt    *gtime.Time `json:"deleted_at"`
}
