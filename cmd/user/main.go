package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mfyai/mfydemo/internal/app/user/api"
	"github.com/mfyai/mfydemo/internal/app/user/config"
	"github.com/mfyai/mfydemo/internal/app/user/handler"
	"github.com/mfyai/mfydemo/internal/app/user/service"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "configs/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx, err := service.NewServiceContext(c)
	if err != nil {
		log.Fatalf("failed to init service context: %v", err)
	}

	server := zrpc.MustNewServer(c.RpcServerConf, func(s *grpc.Server) {
		api.RegisterUserServiceServer(s, handler.NewUserHandler(ctx.UserUsecase))
	})

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	server.Start()
}
