package service

import (
	"context"
	"errors"
	"strings"

	"github.com/mfyai/mfydemo/internal/dao"
	"github.com/mfyai/mfydemo/internal/model/entity"
	"github.com/mfyai/mfydemo/internal/security"
)

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrConflict     = errors.New("resource conflict")
	ErrNotFound     = errors.New("resource not found")
)

type CreateUserInput struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	Status   int
}

type UserService struct {
	dao dao.UserDAO
}

func NewUserService(dao dao.UserDAO) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) Create(ctx context.Context, in CreateUserInput) (*entity.User, error) {
	if s.dao == nil {
		return nil, errors.New("user dao not configured")
	}
	in.Username = strings.TrimSpace(in.Username)
	if in.Username == "" || strings.TrimSpace(in.Password) == "" {
		return nil, ErrInvalidParam
	}
	if in.Status == 0 {
		in.Status = 1
	}

	existing, err := s.dao.GetByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrConflict
	}
	if in.Email != "" {
		existing, err = s.dao.GetByEmail(ctx, in.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrConflict
		}
	}
	if in.Phone != "" {
		existing, err = s.dao.GetByPhone(ctx, in.Phone)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrConflict
		}
	}

	hashed, err := security.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: in.Username,
		Password: hashed,
		Nickname: in.Nickname,
		Email:    in.Email,
		Phone:    in.Phone,
		Status:   in.Status,
	}

	return s.dao.Create(ctx, user)
}

type ListUsersInput struct {
	Page           int
	PageSize       int
	Keyword        string
	Status         *int
	IncludeDeleted bool
}

type ListUsersOutput struct {
	Total int64
	List  []*entity.User
}

func (s *UserService) List(ctx context.Context, in ListUsersInput) (*ListUsersOutput, error) {
	if s.dao == nil {
		return nil, errors.New("user dao not configured")
	}

	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}

	total, list, err := s.dao.List(ctx, dao.UserListQuery{
		Page:           in.Page,
		PageSize:       in.PageSize,
		Keyword:        strings.TrimSpace(in.Keyword),
		Status:         in.Status,
		IncludeDeleted: in.IncludeDeleted,
	})
	if err != nil {
		return nil, err
	}

	return &ListUsersOutput{
		Total: total,
		List:  list,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	if s.dao == nil {
		return nil, errors.New("user dao not configured")
	}
	if id <= 0 {
		return nil, ErrInvalidParam
	}

	user, err := s.dao.GetByID(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user, nil
}
