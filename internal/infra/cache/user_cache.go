package cache

import (
	"context"
	"strconv"
	"time"

	dom "github.com/mfyai/mfydemo/internal/domain/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// NewUserCacheRepo creates a Redis-backed cache repo.
func NewUserCacheRepo(client *redis.Redis, ttl time.Duration) dom.UserCacheRepo {
	return &userCacheRepo{
		client: client,
		ttl:    ttl,
	}
}

type userCacheRepo struct {
	client *redis.Redis
	ttl    time.Duration
}

func (r *userCacheRepo) GetStatus(ctx context.Context, id int64) (int32, bool, error) {
	key := dom.StatusCacheKeyPrefix + strconv.FormatInt(id, 10)
	val, err := r.client.GetCtx(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}
	num, parseErr := strconv.Atoi(val)
	if parseErr != nil {
		return 0, false, parseErr
	}
	return int32(num), true, nil
}

func (r *userCacheRepo) SetStatus(ctx context.Context, id int64, status int32) error {
	key := dom.StatusCacheKeyPrefix + strconv.FormatInt(id, 10)
	return r.client.SetexCtx(ctx, key, strconv.Itoa(int(status)), int(r.ttl/time.Second))
}

func (r *userCacheRepo) DelStatus(ctx context.Context, id int64) error {
	key := dom.StatusCacheKeyPrefix + strconv.FormatInt(id, 10)
	_, err := r.client.DelCtx(ctx, key)
	return err
}
