package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mfyai/mfydemo/internal/domain/entity"
	"github.com/mfyai/mfydemo/internal/domain/repository"
	"github.com/mfyai/mfydemo/pkg/constant"
	"go.uber.org/zap"
)

type UserService struct {
	repo     repository.UserRepository
	cacheTTL time.Duration
	logger   *zap.Logger
}

func NewUserService(repo repository.UserRepository, cacheTTL time.Duration, logger *zap.Logger) *UserService {
	return &UserService{
		repo:     repo,
		cacheTTL: cacheTTL,
		logger:   logger,
	}
}

func (s *UserService) Add(ctx context.Context, u entity.User) (int64, error) {
	if err := s.validateCreate(u); err != nil {
		return 0, err
	}

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.Version = 1

	s.logger.Info("adding user", zap.String("email", u.Email), zap.Int32("status", u.Status))

	id, err := s.repo.Save(ctx, u)
	if err != nil {
		s.logger.Error("save user failed", zap.Error(err))
		return 0, err
	}

	if err := s.repo.SetStatusCache(ctx, id, u.Status, u.Version, s.cacheTTL); err != nil {
		s.logger.Error("set status cache failed", zap.Int64("id", id), zap.Error(err))
		return id, err
	}

	return id, nil
}

func (s *UserService) Update(ctx context.Context, u entity.User) (int64, error) {
	if u.ID <= 0 {
		return 0, fmt.Errorf("%w: user id is required", constant.ErrInvalidArgument)
	}

	if err := s.validateUpdate(u); err != nil {
		return 0, err
	}

	u.UpdatedAt = time.Now()
	s.logger.Info("updating user", zap.Int64("id", u.ID), zap.Int64("version", u.Version), zap.Int32("status", u.Status))

	newVersion, err := s.repo.Update(ctx, u)
	if err != nil {
		s.logger.Error("update user failed", zap.Int64("id", u.ID), zap.Error(err))
		return 0, err
	}

	if err := s.repo.SetStatusCache(ctx, u.ID, u.Status, newVersion, s.cacheTTL); err != nil {
		s.logger.Error("refresh status cache failed", zap.Int64("id", u.ID), zap.Error(err))
		return newVersion, err
	}

	return newVersion, nil
}

func (s *UserService) Remove(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: user id is required", constant.ErrInvalidArgument)
	}

	s.logger.Info("removing user", zap.Int64("id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("delete user failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	if err := s.repo.DelStatusCache(ctx, id); err != nil {
		s.logger.Error("delete status cache failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (s *UserService) Get(ctx context.Context, id int64) (*entity.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: user id is required", constant.ErrInvalidArgument)
	}

	s.logger.Info("fetching user", zap.Int64("id", id))

	cached, err := s.repo.GetStatusCache(ctx, id)
	if err != nil {
		s.logger.Error("get status cache failed", zap.Int64("id", id), zap.Error(err))
	}

	user, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("get user from repo failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if cached == nil || cached.Version != user.Version || cached.Status != user.Status {
		if err := s.repo.SetStatusCache(ctx, user.ID, user.Status, user.Version, s.cacheTTL); err != nil {
			s.logger.Error("backfill status cache failed", zap.Int64("id", user.ID), zap.Error(err))
		}
	}

	return user, nil
}

func (s *UserService) validateCreate(u entity.User) error {
	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("%w: user name is required", constant.ErrInvalidArgument)
	}
	if err := s.validateEmail(u.Email); err != nil {
		return err
	}
	return s.validateStatus(u.Status)
}

func (s *UserService) validateUpdate(u entity.User) error {
	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("%w: user name is required", constant.ErrInvalidArgument)
	}
	if err := s.validateEmail(u.Email); err != nil {
		return err
	}
	if u.Version <= 0 {
		return fmt.Errorf("%w: user version is required", constant.ErrInvalidArgument)
	}
	return s.validateStatus(u.Status)
}

func (s *UserService) validateStatus(status int32) error {
	switch status {
	case constant.UserStatusUnknown, constant.UserStatusActive, constant.UserStatusDisabled:
		return nil
	default:
		return fmt.Errorf("%w: invalid user status %d", constant.ErrInvalidArgument, status)
	}
}

func (s *UserService) validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		return fmt.Errorf("%w: invalid email", constant.ErrInvalidArgument)
	}
	return nil
}
