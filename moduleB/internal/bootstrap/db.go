package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/frame/g"
)

const userTableDDL = `
CREATE TABLE IF NOT EXISTS ` + "`users`" + ` (
  ` + "`id`" + ` BIGINT NOT NULL AUTO_INCREMENT,
  ` + "`username`" + ` VARCHAR(64) NOT NULL,
  ` + "`password`" + ` VARCHAR(255) NOT NULL COMMENT 'hash+salt',
  ` + "`nickname`" + ` VARCHAR(64) DEFAULT NULL,
  ` + "`email`" + ` VARCHAR(128) DEFAULT NULL,
  ` + "`phone`" + ` VARCHAR(32) DEFAULT NULL,
  ` + "`status`" + ` TINYINT NOT NULL DEFAULT 1,
  ` + "`created_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ` + "`updated_at`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  ` + "`deleted_at`" + ` DATETIME DEFAULT NULL,
  PRIMARY KEY (` + "`id`" + `),
  UNIQUE KEY ` + "`uk_users_username`" + ` (` + "`username`" + `),
  UNIQUE KEY ` + "`uk_users_email`" + ` (` + "`email`" + `),
  UNIQUE KEY ` + "`uk_users_phone`" + ` (` + "`phone`" + `),
  KEY ` + "`idx_users_status`" + ` (` + "`status`" + `),
  KEY ` + "`idx_users_deleted_at`" + ` (` + "`deleted_at`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

func EnsureDatabase(ctx context.Context) error {
	cfg := g.Cfg()
	user := cfg.MustGet(ctx, "database.bootstrap.user").String()
	pass := cfg.MustGet(ctx, "database.bootstrap.pass").String()
	host := cfg.MustGet(ctx, "database.bootstrap.host").String()
	port := cfg.MustGet(ctx, "database.bootstrap.port").Int()
	name := cfg.MustGet(ctx, "database.bootstrap.name").String()

	if user == "" || host == "" || port == 0 || name == "" {
		return fmt.Errorf("missing database bootstrap config")
	}

	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port)
	rootDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	if _, err := rootDB.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS `"+name+"`"); err != nil {
		return err
	}

	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
	appDB, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		return err
	}
	defer appDB.Close()

	if _, err := appDB.ExecContext(ctx, userTableDDL); err != nil {
		return err
	}

	return nil
}
