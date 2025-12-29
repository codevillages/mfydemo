package controller

import (
	"context"
	"regexp"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"mfyai/mfydemo/module_user/internal/model"
	"mfyai/mfydemo/module_user/internal/response"
	"mfyai/mfydemo/module_user/internal/service"
)

type UserController struct {
	svc *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		svc: service.NewUserService(),
	}
}

type CreateUserReq struct {
	g.Meta          `path:"/" method:"post" tags:"User" summary:"Create user"`
	Username        string `json:"username" v:"required|length:4,32"`
	Password        string `json:"password" v:"required|length:8,64"`
	PasswordConfirm string `json:"passwordConfirm" v:"required|same:Password"`
	Email           string `json:"email" v:"email"`
	Phone           string `json:"phone"`
	Nickname        string `json:"nickname"`
	Avatar          string `json:"avatar"`
}

type CreateUserRes struct {
	Id uint64 `json:"id"`
}

// Create 用户创建。
func (c *UserController) Create(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	if req.Password != req.PasswordConfirm {
		return nil, gerror.NewCode(response.CodeBadReq, "password confirmation mismatch")
	}
	if req.Phone != "" {
		match, _ := regexp.MatchString(`^\d{11}$`, req.Phone)
		if !match {
			return nil, gerror.NewCode(response.CodeBadReq, "invalid phone format")
		}
	}
	in := &model.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   0,
	}
	id, err := c.svc.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return &CreateUserRes{Id: id}, nil
}

// Healthz 健康检查。
func Healthz(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"code":    response.CodeOK.Code(),
		"message": response.CodeOK.Message(),
		"data":    g.Map{},
	})
}
