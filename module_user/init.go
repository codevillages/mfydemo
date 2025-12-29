package module_user

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"mfyai/mfydemo/module_user/internal/bootstrap"
	"mfyai/mfydemo/module_user/internal/controller"
	"mfyai/mfydemo/module_user/internal/middleware"
)

// Init 注册 module_user 能力：建表、路由与健康检查。
func Init(ctx context.Context, s *ghttp.Server) error {
	if err := bootstrap.EnsureDatabase(ctx); err != nil {
		return err
	}

	group := s.Group("/api/v1")
	ctrl := controller.NewUserController()
	group.Group("/users", func(gp *ghttp.RouterGroup) {
		gp.Bind(ctrl)
	})

	// 健康检查
	s.BindHandler("/healthz", controller.Healthz)
	return nil
}

// 暴露内部中间件给主入口使用。
var (
	RequestIDMiddleware = middleware.RequestID
	AuthMiddleware      = middleware.Auth
	ResponseMiddleware  = middleware.Response
)
