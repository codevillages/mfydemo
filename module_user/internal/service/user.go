package service

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"

	"mfyai/mfydemo/module_user/internal/dao/mysql"
	"mfyai/mfydemo/module_user/internal/model"
	"mfyai/mfydemo/module_user/internal/response"
)

type UserService struct {
	dao *mysql.UserDAO
}

func NewUserService() *UserService {
	return &UserService{
		dao: mysql.NewUserDAO(),
	}
}

// Create handles user creation with unique validation and password hashing.
func (s *UserService) Create(ctx context.Context, in *model.CreateUserInput) (uint64, error) {
	conflictField, err := s.dao.FindConflict(ctx, in.Username, in.Email, in.Phone)
	if err != nil {
		return 0, err
	}
	if conflictField != "" {
		return 0, gerror.NewCode(response.CodeConflict, fmt.Sprintf("%s already exists", conflictField))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil {
		return 0, err
	}
	id, err := s.dao.Create(ctx, in, string(hash))
	if err != nil {
		return 0, err
	}
	return id, nil
}
