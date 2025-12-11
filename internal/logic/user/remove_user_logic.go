package user

import (
	"context"

	"github.com/mfyai/mfydemo/internal/svc"
	userpb "github.com/mfyai/mfydemo/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RemoveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveUserLogic {
	return &RemoveUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveUserLogic) RemoveUser(in *userpb.RemoveUserReq) (*emptypb.Empty, error) {
	logger := l.svcCtx.Logger.With(zap.String("method", "RemoveUser"), zap.Int64("id", in.GetId()))
	logger.Info("incoming remove user", zap.Int64("id", in.GetId()))

	if err := l.svcCtx.UserService.Remove(l.ctx, in.GetId()); err != nil {
		logger.Error("remove user failed", zap.Error(err))
		return nil, mapError(err)
	}

	logger.Info("remove user success", zap.Int64("id", in.GetId()))
	return &emptypb.Empty{}, nil
}
