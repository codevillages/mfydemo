package usecase

import (
	"context"

	dom "github.com/mfyai/mfydemo/internal/domain/user"
	"go.uber.org/zap"
)

// UserUsecase orchestrates user operations.
type UserUsecase struct {
	repo           dom.UserRepo
	cache          dom.UserCacheRepo
	log            *zap.Logger
	hashPassword   dom.HashPassword
	compareHash    dom.CompareHash
}

// NewUserUsecase constructs a usecase instance.
func NewUserUsecase(repo dom.UserRepo, cache dom.UserCacheRepo, log *zap.Logger, hash dom.HashPassword, compare dom.CompareHash) *UserUsecase {
	return &UserUsecase{
		repo:         repo,
		cache:        cache,
		log:          log,
		hashPassword: hash,
		compareHash:  compare,
	}
}

// AddUser creates a new user.
func (u *UserUsecase) AddUser(ctx context.Context, username, password, email string, status int32) (int64, error) {
	u.log.Info("AddUser request received", zap.String("username", username))

	if err := dom.ValidateNew(username, password, email, status); err != nil {
		return 0, err
	}

	if existing, err := u.repo.FindByUsername(ctx, username); err == nil && existing != nil {
		return 0, dom.ErrUserExists
	}

	hashed, err := u.hashPassword(password)
	if err != nil {
		u.log.Error("hash password failed", zap.Error(err))
		return 0, dom.ErrHashPassword
	}

	entity := &dom.User{
		Username: username,
		Password: hashed,
		Email:    email,
		Status:   status,
	}
	id, err := u.repo.Create(ctx, entity)
	if err != nil {
		u.log.Error("db insert failed", zap.Error(err))
		return 0, err
	}

	if err := u.cache.SetStatus(ctx, id, status); err != nil {
		u.log.Error("cache set status failed", zap.Error(err), zap.Int64("user_id", id))
	}

	u.log.Info("AddUser success", zap.Int64("user_id", id))
	return id, nil
}

// UpdateUser updates user fields.
func (u *UserUsecase) UpdateUser(ctx context.Context, id int64, password, email string, status int32) error {
	u.log.Info("UpdateUser request received", zap.Int64("user_id", id))

	if err := dom.ValidateUpdate(password, email, status); err != nil {
		return err
	}

	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	hashed := existing.Password
	if password != "" {
		hashed, err = u.hashPassword(password)
		if err != nil {
			u.log.Error("hash password failed", zap.Error(err), zap.Int64("user_id", id))
			return dom.ErrHashPassword
		}
	}

	entity := &dom.User{
		ID:       id,
		Username: existing.Username,
		Password: hashed,
		Email:    email,
		Status:   status,
	}
	if entity.Email == "" {
		entity.Email = existing.Email
	}

	if err := u.repo.Update(ctx, entity); err != nil {
		u.log.Error("db update failed", zap.Error(err), zap.Int64("user_id", id))
		return err
	}

	if err := u.cache.SetStatus(ctx, id, status); err != nil {
		u.log.Error("cache set status failed", zap.Error(err), zap.Int64("user_id", id))
	}

	u.log.Info("UpdateUser success", zap.Int64("user_id", id))
	return nil
}

// RemoveUser deletes a user (hard delete).
func (u *UserUsecase) RemoveUser(ctx context.Context, id int64) error {
	u.log.Info("RemoveUser request received", zap.Int64("user_id", id))

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := u.cache.DelStatus(ctx, id); err != nil {
		u.log.Error("cache delete status failed", zap.Error(err), zap.Int64("user_id", id))
	}

	u.log.Info("RemoveUser success", zap.Int64("user_id", id))
	return nil
}

// GetUser retrieves user info with cached status.
func (u *UserUsecase) GetUser(ctx context.Context, id int64) (*dom.User, error) {
	u.log.Info("GetUser request received", zap.Int64("user_id", id))

	status, ok, err := u.cache.GetStatus(ctx, id)
	if err != nil {
		u.log.Error("cache get status failed", zap.Error(err), zap.Int64("user_id", id))
	}

	entity, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok {
		entity.Status = status
	} else {
		if err := u.cache.SetStatus(ctx, id, entity.Status); err != nil {
			u.log.Error("cache backfill status failed", zap.Error(err), zap.Int64("user_id", id))
		}
	}

	u.log.Info("GetUser success", zap.Int64("user_id", id))
	return entity, nil
}
