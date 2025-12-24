package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/mfyai/mfydemo/internal/dao"
	"github.com/mfyai/mfydemo/internal/model/entity"
)

var errNotImplemented = errors.New("not implemented")

const userTable = "users"

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (d *UserDAO) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}
	result, err := d.model(ctx).Data(g.Map{
		"username": user.Username,
		"password": user.Password,
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
		"status":   user.Status,
	}).Insert()
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return d.GetByID(ctx, id, true)
}

func (d *UserDAO) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return d.getOne(ctx, g.Map{"username": username})
}

func (d *UserDAO) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return d.getOne(ctx, g.Map{"email": email})
}

func (d *UserDAO) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	return d.getOne(ctx, g.Map{"phone": phone})
}

func (d *UserDAO) GetByID(ctx context.Context, id int64, includeDeleted bool) (*entity.User, error) {
	model := d.model(ctx).Where("id", id)
	if !includeDeleted {
		model = model.Where("deleted_at", nil)
	}
	return d.getOneWithModel(ctx, model)
}

func (d *UserDAO) List(ctx context.Context, query dao.UserListQuery) (int64, []*entity.User, error) {
	model := d.model(ctx)
	if !query.IncludeDeleted {
		model = model.Where("deleted_at", nil)
	}
	if query.Status != nil {
		model = model.Where("status", *query.Status)
	}
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		model = model.WhereOrLike("username", like).
			WhereOrLike("nickname", like).
			WhereOrLike("email", like).
			WhereOrLike("phone", like)
	}

	total, err := model.Clone().Count()
	if err != nil {
		return 0, nil, err
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	offset := (query.Page - 1) * query.PageSize

	var users []*entity.User
	if err := model.Limit(query.PageSize).Offset(offset).Order("id desc").Scan(&users); err != nil {
		return 0, nil, err
	}

	return total, users, nil
}

func (d *UserDAO) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	return nil, errNotImplemented
}

func (d *UserDAO) Delete(ctx context.Context, id int64) error {
	return errNotImplemented
}

func (d *UserDAO) HardDelete(ctx context.Context, id int64) error {
	return errNotImplemented
}

func (d *UserDAO) model(ctx context.Context) *gdb.Model {
	return g.DB().Model(userTable).Ctx(ctx)
}

func (d *UserDAO) getOne(ctx context.Context, where g.Map) (*entity.User, error) {
	model := d.model(ctx).Where(where).Where("deleted_at", nil)
	return d.getOneWithModel(ctx, model)
}

func (d *UserDAO) getOneWithModel(ctx context.Context, model *gdb.Model) (*entity.User, error) {
	record, err := model.One()
	if err != nil {
		return nil, err
	}
	if record == nil || record.IsEmpty() {
		return nil, nil
	}
	var user entity.User
	if err := record.Struct(&user); err != nil {
		return nil, err
	}
	if user.DeletedAt != nil && user.DeletedAt.After(time.Time{}) {
		return nil, nil
	}
	return &user, nil
}
