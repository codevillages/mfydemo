package repository

import (
	"context"
	"time"

	"github.com/mfyai/mfydemo/internal/domain/entity"
)

type StatusCache struct {
	Status  int32
	Version int64
}

type UserRepository interface {
	Save(ctx context.Context, u entity.User) (int64, error)
	Update(ctx context.Context, u entity.User) (int64, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*entity.User, error)
	SetStatusCache(ctx context.Context, id int64, status int32, version int64, ttl time.Duration) error
	GetStatusCache(ctx context.Context, id int64) (*StatusCache, error)
	DelStatusCache(ctx context.Context, id int64) error
}
