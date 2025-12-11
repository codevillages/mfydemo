package user

import (
	"context"

	"github.com/mfyai/mfydemo/internal/domain/entity"
	"github.com/mfyai/mfydemo/internal/svc"
	userpb "github.com/mfyai/mfydemo/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(in *userpb.UpdateUserReq) (*emptypb.Empty, error) {
	logger := l.svcCtx.Logger.With(zap.String("method", "UpdateUser"), zap.Int64("id", in.GetId()), zap.Int64("version", in.GetVersion()))
	logger.Info("incoming update user", zap.Int64("id", in.GetId()), zap.Int32("status", in.GetStatus()))

	_, err := l.svcCtx.UserService.Update(l.ctx, entity.User{
		ID:      in.GetId(),
		Name:    in.GetName(),
		Email:   in.GetEmail(),
		Status:  in.GetStatus(),
		Version: in.GetVersion(),
	})
	if err != nil {
		logger.Error("update user failed", zap.Error(err))
		return nil, mapError(err)
	}

	logger.Info("update user success", zap.Int64("id", in.GetId()))
	return &emptypb.Empty{}, nil
}
