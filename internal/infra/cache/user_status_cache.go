package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mfyai/mfydemo/internal/domain/repository"
	"github.com/mfyai/mfydemo/pkg/constant"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.uber.org/zap"
)

type UserStatusCache struct {
	client *redis.Redis
	logger *zap.Logger
}

type statusCacheValue struct {
	Status  int32 `json:"status"`
	Version int64 `json:"version"`
}

func NewUserStatusCache(client *redis.Redis, logger *zap.Logger) *UserStatusCache {
	return &UserStatusCache{
		client: client,
		logger: logger,
	}
}

func (c *UserStatusCache) Set(ctx context.Context, id int64, status int32, version int64, ttl time.Duration) error {
	raw, err := json.Marshal(statusCacheValue{
		Status:  status,
		Version: version,
	})
	if err != nil {
		return err
	}

	return c.client.SetexCtx(ctx, constant.StatusCacheKey(id), string(raw), int(ttl/time.Second))
}

func (c *UserStatusCache) Get(ctx context.Context, id int64) (*repository.StatusCache, error) {
	val, err := c.client.GetCtx(ctx, constant.StatusCacheKey(id))
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var payload statusCacheValue
	if err := json.Unmarshal([]byte(val), &payload); err != nil {
		c.logger.Error("unmarshal status cache failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return &repository.StatusCache{
		Status:  payload.Status,
		Version: payload.Version,
	}, nil
}

func (c *UserStatusCache) Del(ctx context.Context, id int64) error {
	_, err := c.client.DelCtx(ctx, constant.StatusCacheKey(id))
	return err
}
