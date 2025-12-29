package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"mfyai/mfydemo/module_user/internal/model"
	"mfyai/mfydemo/module_user/internal/model/entity"
)

type UserDAO struct {
	table  string
	schema string
}

func NewUserDAO() *UserDAO {
	return &UserDAO{
		table:  "users",
		schema: schemaName(),
	}
}

func (d *UserDAO) model() *gdb.Model {
	if d.schema != "" {
		return g.DB().Schema(d.schema).Model(d.table)
	}
	return g.DB().Model(d.table)
}

// Create inserts a user and returns id.
func (d *UserDAO) Create(ctx context.Context, in *model.CreateUserInput, passwordHash string) (uint64, error) {
	table := d.table
	if d.schema != "" {
		table = fmt.Sprintf("`%s`.`%s`", d.schema, d.table)
	} else {
		table = fmt.Sprintf("`%s`", d.table)
	}
	sqlStr := fmt.Sprintf(
		"INSERT INTO %s (username, email, phone, nickname, avatar, password_hash, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		table,
	)
	res, err := g.DB().Exec(ctx, sqlStr,
		in.Username,
		emptyToNil(in.Email),
		emptyToNil(in.Phone),
		emptyToNil(in.Nickname),
		emptyToNil(in.Avatar),
		passwordHash,
		in.Status,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

// FindConflict checks unique fields; returns conflict field name.
func (d *UserDAO) FindConflict(ctx context.Context, username, email, phone string) (string, error) {
	builder := d.model().Where("deleted_at IS NULL")

	conds := make([]gdb.WhereHolder, 0, 3)
	if username != "" {
		conds = append(conds, gdb.WhereHolder{
			Where: "username = ?",
			Args:  []interface{}{username},
		})
	}
	if email != "" {
		conds = append(conds, gdb.WhereHolder{
			Where: "email = ?",
			Args:  []interface{}{email},
		})
	}
	if phone != "" {
		conds = append(conds, gdb.WhereHolder{
			Where: "phone = ?",
			Args:  []interface{}{phone},
		})
	}
	if len(conds) == 0 {
		return "", nil
	}
	builder = builder.Where(conds[0])
	for i := 1; i < len(conds); i++ {
		builder = builder.WhereOr(conds[i])
	}
	var u entity.User
	err := builder.Fields("id, username, email, phone").Limit(1).Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	if u.Username == username && username != "" {
		return "username", nil
	}
	if u.Email == email && email != "" {
		return "email", nil
	}
	if u.Phone == phone && phone != "" {
		return "phone", nil
	}
	return "", nil
}

func emptyToNil(v string) interface{} {
	if v == "" {
		return nil
	}
	return v
}

func schemaName() string {
	cfg := g.DB().GetConfig()
	if cfg != nil && cfg.Name != "" {
		return cfg.Name
	}
	return "mfydemo_user"
}
