package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type User struct {
	Id         int64     `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	Status     int32     `db:"status"`
	Version    int64     `db:"version"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type UserModel interface {
	Insert(ctx context.Context, data *User) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, session sqlx.Session, data *User, previousVersion int64) (sql.Result, error)
	Delete(ctx context.Context, session sqlx.Session, id int64) (sql.Result, error)
	TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
}

type defaultUserModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &defaultUserModel{
		conn:  conn,
		table: "`users`",
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`name`, `email`, `status`, `version`, `create_time`, `update_time`) values (?, ?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.Name, data.Email, data.Status, data.Version, data.CreateTime, data.UpdateTime)
}

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	query := fmt.Sprintf("select `id`, `name`, `email`, `status`, `version`, `create_time`, `update_time` from %s where `id` = ? limit 1", m.table)

	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (m *defaultUserModel) Update(ctx context.Context, session sqlx.Session, data *User, previousVersion int64) (sql.Result, error) {
	query := fmt.Sprintf("update %s set `name` = ?, `email` = ?, `status` = ?, `version` = ?, `update_time` = ? where `id` = ? and `version` = ?", m.table)

	execFn := m.conn.ExecCtx
	if session != nil {
		execFn = session.ExecCtx
	}

	return execFn(ctx, query, data.Name, data.Email, data.Status, data.Version, data.UpdateTime, data.Id, previousVersion)
}

func (m *defaultUserModel) Delete(ctx context.Context, session sqlx.Session, id int64) (sql.Result, error) {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)

	execFn := m.conn.ExecCtx
	if session != nil {
		execFn = session.ExecCtx
	}

	return execFn(ctx, query, id)
}

func (m *defaultUserModel) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}
