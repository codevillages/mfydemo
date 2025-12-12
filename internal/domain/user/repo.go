package user

import "context"

// UserRepo defines persistent storage operations.
type UserRepo interface {
	Create(ctx context.Context, u *User) (int64, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}

// UserCacheRepo defines cache operations for user status.
type UserCacheRepo interface {
	GetStatus(ctx context.Context, id int64) (int32, bool, error)
	SetStatus(ctx context.Context, id int64, status int32) error
	DelStatus(ctx context.Context, id int64) error
}
