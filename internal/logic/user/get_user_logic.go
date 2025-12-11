package user

import (
	"context"

	"github.com/mfyai/mfydemo/internal/domain/entity"
	"github.com/mfyai/mfydemo/internal/svc"
	userpb "github.com/mfyai/mfydemo/proto"
	"go.uber.org/zap"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(in *userpb.GetUserReq) (*userpb.GetUserResp, error) {
	logger := l.svcCtx.Logger.With(zap.String("method", "GetUser"), zap.Int64("id", in.GetId()))
	logger.Info("incoming get user", zap.Int64("id", in.GetId()))

	userEntity, err := l.svcCtx.UserService.Get(l.ctx, in.GetId())
	if err != nil {
		logger.Error("get user failed", zap.Error(err))
		return nil, mapError(err)
	}

	return &userpb.GetUserResp{
		User: convertToProto(userEntity),
	}, nil
}

func convertToProto(u *entity.User) *userpb.User {
	if u == nil {
		return nil
	}

	return &userpb.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Status:    u.Status,
		Version:   u.Version,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}
