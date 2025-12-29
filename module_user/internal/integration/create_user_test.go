package integration

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"mfyai/mfydemo/module_user/internal/bootstrap"
	"mfyai/mfydemo/module_user/internal/model"
	"mfyai/mfydemo/module_user/internal/service"
)

// TestCreateUser inserts one user to verify DB connectivity and create flow.
func TestCreateUser(t *testing.T) {
	if os.Getenv("GF_GCFG_FILE") == "" {
		t.Skip("GF_GCFG_FILE not set; skip DB integration")
	}

	ctx := context.Background()

	if err := bootstrap.EnsureDatabase(ctx); err != nil {
		t.Fatalf("ensure database: %v", err)
	}

	svc := service.NewUserService()
	username := "testuser_" + time.Now().Format("150405")
	id, err := svc.Create(ctx, &model.CreateUserInput{
		Username: username,
		Password: "Secret123!",
		Email:    username + "@example.com",
		Phone:    "13800000000",
		Nickname: "Tester",
		Avatar:   "",
		Status:   0,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero id")
	}
	t.Logf("created user id=%d username=%s", id, username)
}
