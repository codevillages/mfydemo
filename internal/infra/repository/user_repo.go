package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mfyai/mfydemo/internal/domain/entity"
	domainrepo "github.com/mfyai/mfydemo/internal/domain/repository"
	"github.com/mfyai/mfydemo/internal/infra/cache"
	"github.com/mfyai/mfydemo/internal/infra/dao"
	"github.com/mfyai/mfydemo/pkg/constant"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
)

type userRepository struct {
	model dao.UserModel
	cache *cache.UserStatusCache
	log   *zap.Logger
}

func NewUserRepository(model dao.UserModel, cache *cache.UserStatusCache, log *zap.Logger) domainrepo.UserRepository {
	return &userRepository{
		model: model,
		cache: cache,
		log:   log,
	}
}

func (r *userRepository) Save(ctx context.Context, u entity.User) (int64, error) {
	data := dao.User{
		Name:       u.Name,
		Email:      u.Email,
		Status:     u.Status,
		Version:    u.Version,
		CreateTime: u.CreatedAt,
		UpdateTime: u.UpdatedAt,
	}

	res, err := r.model.Insert(ctx, &data)
	if err != nil {
		r.log.Error("insert user failed", zap.String("email", u.Email), zap.Error(err))
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.log.Error("get last insert id failed", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (r *userRepository) Update(ctx context.Context, u entity.User) (int64, error) {
	newVersion := u.Version + 1
	dbUser := dao.User{
		Id:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Status:     u.Status,
		Version:    newVersion,
		UpdateTime: u.UpdatedAt,
	}

	err := r.model.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		res, err := r.model.Update(ctx, session, &dbUser, u.Version)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			exists, findErr := r.model.FindOne(ctx, u.ID)
			if findErr != nil {
				if findErr == sqlx.ErrNotFound || findErr == sql.ErrNoRows {
					return constant.ErrUserNotFound
				}
				return findErr
			}
			r.log.Warn("version conflict on update", zap.Int64("id", u.ID), zap.Int64("version", u.Version), zap.Any("existing", exists.Version))
			return constant.ErrVersionConflict
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return newVersion, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	err := r.model.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		res, err := r.model.Delete(ctx, session, id)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			return constant.ErrUserNotFound
		}

		return nil
	})
	if err != nil {
		r.log.Error("delete user failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (r *userRepository) Get(ctx context.Context, id int64) (*entity.User, error) {
	resp, err := r.model.FindOne(ctx, id)
	if err != nil {
		if err == sqlx.ErrNotFound || err == sql.ErrNoRows {
			return nil, constant.ErrUserNotFound
		}
		return nil, err
	}

	return &entity.User{
		ID:        resp.Id,
		Name:      resp.Name,
		Email:     resp.Email,
		Status:    resp.Status,
		Version:   resp.Version,
		CreatedAt: resp.CreateTime,
		UpdatedAt: resp.UpdateTime,
	}, nil
}

func (r *userRepository) SetStatusCache(ctx context.Context, id int64, status int32, version int64, ttlDuration time.Duration) error {
	return r.cache.Set(ctx, id, status, version, ttlDuration)
}

func (r *userRepository) GetStatusCache(ctx context.Context, id int64) (*domainrepo.StatusCache, error) {
	return r.cache.Get(ctx, id)
}

func (r *userRepository) DelStatusCache(ctx context.Context, id int64) error {
	return r.cache.Del(ctx, id)
}
