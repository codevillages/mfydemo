package main

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"

	"mfyai/mfydemo/module_user"
)

func main() {
	// 配置目录，支持通过 gf.gcfg.path 覆盖
	if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		adapter.SetPath("config")
	}

	s := g.Server()
	registerMiddlewares(s)

	ctx := context.Background()
	if err := module_user.Init(ctx, s); err != nil {
		g.Log().Fatal(ctx, err)
	}

	s.Run()
}

func registerMiddlewares(s *ghttp.Server) {
	s.Use(
		module_user.RequestIDMiddleware,
		ghttp.MiddlewareAccessLog,
		module_user.ResponseMiddleware,
		module_user.AuthMiddleware,
	)
}
