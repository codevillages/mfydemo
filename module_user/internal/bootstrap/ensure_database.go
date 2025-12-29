package bootstrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EnsureDatabase 创建 database/schema 及 users 表（若不存在）。
func EnsureDatabase(ctx context.Context) error {
	db := g.DB()
	if db == nil {
		return fmt.Errorf("database not configured")
	}

	schema := dbName(ctx)
	if schema == "" {
		return fmt.Errorf("database name not configured")
	}

	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", schema)
	if _, err := db.Exec(ctx, createDB); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	ddl := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s.users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(128) NULL UNIQUE,
    phone VARCHAR(16) NULL UNIQUE,
    nickname VARCHAR(64) NULL,
    avatar VARCHAR(255) NULL,
    password_hash VARCHAR(255) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY (id),
    INDEX idx_status_created_at (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`, schema)
	if _, err := db.Exec(ctx, ddl); err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	return nil
}

// WithTx is a helper to run a function in transaction.
func WithTx(ctx context.Context, fn func(ctx context.Context, tx gdb.TX) error) error {
	return g.DB().Transaction(ctx, fn)
}

// dbName reads target schema name from config or db config.
func dbName(ctx context.Context) string {
	cfg := g.DB().GetConfig()
	if cfg != nil && strings.TrimSpace(cfg.Name) != "" {
		return strings.TrimSpace(cfg.Name)
	}
	name := g.Cfg().MustGet(ctx, "database.default.name").String()
	if name != "" {
		return name
	}
	return "mfydemo_user"
}
