package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/mfyai/mfydemo/internal/bootstrap"
	"github.com/mfyai/mfydemo/internal/controller"
	"github.com/mfyai/mfydemo/internal/dao/mysql"
	"github.com/mfyai/mfydemo/internal/middleware"
	"github.com/mfyai/mfydemo/internal/service"
)

type createUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"data"`
}

func TestCreateUserIntegration(t *testing.T) {
	ctx := context.Background()
	if err := bootstrap.EnsureDatabase(ctx); err != nil {
		t.Fatalf("bootstrap database failed: %v", err)
	}

	s := g.Server()
	s.SetPort(18080)
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Auth)
		userDAO := mysql.NewUserDAO()
		userService := service.NewUserService(userDAO)
		userController := controller.NewUserController(userService)
		group.POST("/users", userController.Create)
	})

	if err := s.Start(); err != nil {
		t.Fatalf("start server failed: %v", err)
	}
	defer s.Shutdown()

	username := fmt.Sprintf("itest_%d", time.Now().UnixNano())
	payload := map[string]interface{}{
		"username": username,
		"password": "testpass123",
		"nickname": "itest",
		"status":   1,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload failed: %v", err)
	}

	resp, err := http.Post("http://127.0.0.1:18080/api/v1/users", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http 200, got %d", resp.StatusCode)
	}

	var out createUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	if out.Code != 0 {
		t.Fatalf("expected code 0, got %d", out.Code)
	}
	if out.Data.Username != username {
		t.Fatalf("expected username %s, got %s", username, out.Data.Username)
	}
}
