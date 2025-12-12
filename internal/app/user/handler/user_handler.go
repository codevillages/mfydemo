package handler

import (
	"context"
	"time"

	"github.com/mfyai/mfydemo/internal/app/user/api"
	"github.com/mfyai/mfydemo/internal/app/user/usecase"
	dom "github.com/mfyai/mfydemo/internal/domain/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler implements gRPC handlers.
type UserHandler struct {
	api.UnimplementedUserServiceServer
	uc *usecase.UserUsecase
}

// NewUserHandler constructs handler.
func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) AddUser(ctx context.Context, req *api.AddUserRequest) (*api.AddUserResponse, error) {
	id, err := h.uc.AddUser(ctx, req.Username, req.Password, req.Email, req.Status)
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &api.AddUserResponse{Id: id}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (*api.UpdateUserResponse, error) {
	if err := h.uc.UpdateUser(ctx, req.Id, req.Password, req.Email, req.Status); err != nil {
		return nil, toStatusErr(err)
	}
	return &api.UpdateUserResponse{Success: true}, nil
}

func (h *UserHandler) RemoveUser(ctx context.Context, req *api.RemoveUserRequest) (*api.RemoveUserResponse, error) {
	if err := h.uc.RemoveUser(ctx, req.Id); err != nil {
		return nil, toStatusErr(err)
	}
	return &api.RemoveUserResponse{Success: true}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	entity, err := h.uc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &api.GetUserResponse{
		Id:        entity.ID,
		Username:  entity.Username,
		Email:     entity.Email,
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: entity.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func toStatusErr(err error) error {
	switch err {
	case nil:
		return nil
	case dom.ErrInvalidAccount, dom.ErrInvalidEmail, dom.ErrInvalidStatus:
		return status.Error(codes.InvalidArgument, err.Error())
	case dom.ErrUserExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case dom.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case dom.ErrHashPassword:
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
