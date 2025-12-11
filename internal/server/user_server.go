package server

import (
	"context"

	userlogic "github.com/mfyai/mfydemo/internal/logic/user"
	"github.com/mfyai/mfydemo/internal/svc"
	userpb "github.com/mfyai/mfydemo/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	svcCtx *svc.ServiceContext
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) AddUser(ctx context.Context, in *userpb.AddUserReq) (*userpb.AddUserResp, error) {
	return userlogic.NewAddUserLogic(ctx, s.svcCtx).AddUser(in)
}

func (s *UserServer) UpdateUser(ctx context.Context, in *userpb.UpdateUserReq) (*emptypb.Empty, error) {
	return userlogic.NewUpdateUserLogic(ctx, s.svcCtx).UpdateUser(in)
}

func (s *UserServer) RemoveUser(ctx context.Context, in *userpb.RemoveUserReq) (*emptypb.Empty, error) {
	return userlogic.NewRemoveUserLogic(ctx, s.svcCtx).RemoveUser(in)
}

func (s *UserServer) GetUser(ctx context.Context, in *userpb.GetUserReq) (*userpb.GetUserResp, error) {
	return userlogic.NewGetUserLogic(ctx, s.svcCtx).GetUser(in)
}
