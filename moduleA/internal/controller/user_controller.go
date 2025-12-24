package controller

import (
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/mfyai/mfydemo/internal/response"
	"github.com/mfyai/mfydemo/internal/service"
)

type UserController struct {
	service *service.UserService
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}

type ListUsersRequest struct {
	Page           int    `json:"page" v:"min:1#page must be >= 1"`
	PageSize       int    `json:"page_size" v:"min:1#page_size must be >= 1"`
	Keyword        string `json:"keyword"`
	Status         *int   `json:"status"`
	IncludeDeleted bool   `json:"include_deleted"`
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Create(r *ghttp.Request) {
	var req CreateUserRequest
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatusExit(http.StatusBadRequest, response.Error(response.CodeInvalidParam, "invalid request"))
		return
	}

	user, err := c.service.Create(r.Context(), service.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidParam:
			r.Response.WriteStatusExit(http.StatusBadRequest, response.Error(response.CodeInvalidParam, "invalid parameter"))
		case service.ErrConflict:
			r.Response.WriteStatusExit(http.StatusConflict, response.Error(response.CodeConflict, "resource conflict"))
		default:
			r.Response.WriteStatusExit(http.StatusInternalServerError, response.Error(response.CodeInternal, "internal error"))
		}
		return
	}

	r.Response.WriteStatusExit(http.StatusOK, response.Success(user))
}

func (c *UserController) List(r *ghttp.Request) {
	var req ListUsersRequest
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatusExit(http.StatusBadRequest, response.Error(response.CodeInvalidParam, "invalid request"))
		return
	}

	result, err := c.service.List(r.Context(), service.ListUsersInput{
		Page:           req.Page,
		PageSize:       req.PageSize,
		Keyword:        req.Keyword,
		Status:         req.Status,
		IncludeDeleted: req.IncludeDeleted,
	})
	if err != nil {
		r.Response.WriteStatusExit(http.StatusInternalServerError, response.Error(response.CodeInternal, "internal error"))
		return
	}

	r.Response.WriteStatusExit(http.StatusOK, response.Success(result))
}

func (c *UserController) Detail(r *ghttp.Request) {
	id := r.Get("id").Int64()
	if id <= 0 {
		r.Response.WriteStatusExit(http.StatusBadRequest, response.Error(response.CodeInvalidParam, "invalid parameter"))
		return
	}

	user, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case service.ErrInvalidParam:
			r.Response.WriteStatusExit(http.StatusBadRequest, response.Error(response.CodeInvalidParam, "invalid parameter"))
		case service.ErrNotFound:
			r.Response.WriteStatusExit(http.StatusNotFound, response.Error(response.CodeNotFound, "user not found"))
		default:
			r.Response.WriteStatusExit(http.StatusInternalServerError, response.Error(response.CodeInternal, "internal error"))
		}
		return
	}

	r.Response.WriteStatusExit(http.StatusOK, response.Success(user))
}
