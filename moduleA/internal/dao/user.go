package dao

import (
	"context"

	"github.com/mfyai/mfydemo/internal/model/entity"
)

type UserDAO interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByID(ctx context.Context, id int64, includeDeleted bool) (*entity.User, error)
	List(ctx context.Context, query UserListQuery) (int64, []*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
	HardDelete(ctx context.Context, id int64) error
}

type UserListQuery struct {
	Page           int
	PageSize       int
	Keyword        string
	Status         *int
	IncludeDeleted bool
}
