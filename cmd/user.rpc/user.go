package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mfyai/mfydemo/internal/config"
	"github.com/mfyai/mfydemo/internal/server"
	"github.com/mfyai/mfydemo/internal/svc"
	userpb "github.com/mfyai/mfydemo/proto"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "cmd/user.rpc/etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcCtx, err := svc.NewServiceContext(c)
	if err != nil {
		log.Fatalf("failed to init service context: %v", err)
	}
	defer svcCtx.SyncLogger()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		userpb.RegisterUserServiceServer(grpcServer, server.NewUserServer(svcCtx))
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
