package user

import (
	"context"

	"github.com/mfyai/mfydemo/internal/domain/entity"
	"github.com/mfyai/mfydemo/internal/svc"
	userpb "github.com/mfyai/mfydemo/proto"
	"go.uber.org/zap"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUserLogic) AddUser(in *userpb.AddUserReq) (*userpb.AddUserResp, error) {
	logger := l.svcCtx.Logger.With(zap.String("method", "AddUser"))
	logger.Info("incoming add user", zap.String("name", in.GetName()), zap.String("email", in.GetEmail()), zap.Int32("status", in.GetStatus()))

	id, err := l.svcCtx.UserService.Add(l.ctx, entity.User{
		Name:   in.GetName(),
		Email:  in.GetEmail(),
		Status: in.GetStatus(),
	})
	if err != nil {
		logger.Error("add user failed", zap.Error(err))
		return nil, mapError(err)
	}

	logger.Info("add user success", zap.Int64("id", id))
	return &userpb.AddUserResp{Id: id}, nil
}
