package persistence

import (
	"context"
	"database/sql"

	dom "github.com/mfyai/mfydemo/internal/domain/user"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// NewUserRepo creates a DB-backed user repository.
func NewUserRepo(conn sqlx.SqlConn) dom.UserRepo {
	return &userRepo{
		conn:  conn,
		table: "users",
	}
}

type userRepo struct {
	conn  sqlx.SqlConn
	table string
}

func (r *userRepo) Create(ctx context.Context, u *dom.User) (int64, error) {
	query := "INSERT INTO " + r.table + " (username, password, email, status, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	res, err := r.conn.ExecCtx(ctx, query, u.Username, u.Password, u.Email, u.Status)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepo) Update(ctx context.Context, u *dom.User) error {
	query := "UPDATE " + r.table + " SET password = ?, email = ?, status = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.conn.ExecCtx(ctx, query, u.Password, u.Email, u.Status, u.ID)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM " + r.table + " WHERE id = ?"
	_, err := r.conn.ExecCtx(ctx, query, id)
	return err
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*dom.User, error) {
	query := "SELECT id, username, password, email, status, created_at, updated_at FROM " + r.table + " WHERE id = ? LIMIT 1"
	var u dom.User
	if err := r.conn.QueryRowCtx(ctx, &u, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, dom.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*dom.User, error) {
	query := "SELECT id, username, password, email, status, created_at, updated_at FROM " + r.table + " WHERE username = ? LIMIT 1"
	var u dom.User
	if err := r.conn.QueryRowCtx(ctx, &u, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, dom.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
