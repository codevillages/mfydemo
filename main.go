package main

import (
	"context"
	"log"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/mfyai/mfydemo/internal/bootstrap"
	"github.com/mfyai/mfydemo/internal/controller"
	"github.com/mfyai/mfydemo/internal/dao/mysql"
	"github.com/mfyai/mfydemo/internal/middleware"
	"github.com/mfyai/mfydemo/internal/service"
)

func main() {
	s := g.Server()
	s.SetPort(8080)

	if err := bootstrap.EnsureDatabase(context.Background()); err != nil {
		log.Fatalf("database bootstrap failed: %v", err)
	}

	userDAO := mysql.NewUserDAO()
	userService := service.NewUserService(userDAO)
	userController := controller.NewUserController(userService)

	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Auth)
		group.POST("/users", userController.Create)
		group.GET("/users", userController.List)
		group.GET("/users/{id}", userController.Detail)
	})

	s.Run()
}
