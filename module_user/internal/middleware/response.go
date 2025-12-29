package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"mfyai/mfydemo/module_user/internal/response"
)

// Response 将处理结果统一包装为 {code,message,data} 并按业务码设置 HTTP 状态码。
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.Written() {
		return
	}

	if err := r.GetError(); err != nil {
		code := gerror.Code(err)
		if code == gcode.CodeNil {
			code = response.CodeInternal
		} else if code.Code() == gcode.CodeValidationFailed.Code() {
			code = response.CodeBadReq
		}
		status := response.HTTPStatus(code.Code())
		r.Response.WriteStatus(status)
		r.Response.WriteJsonExit(g.Map{
			"code":    code.Code(),
			"message": err.Error(),
			"data":    g.Map{},
		})
		return
	}

	res := r.GetHandlerResponse()
	if res == nil {
		res = g.Map{}
	}
	r.Response.WriteStatus(http.StatusOK)
	r.Response.WriteJsonExit(g.Map{
		"code":    response.CodeOK.Code(),
		"message": response.CodeOK.Message(),
		"data":    res,
	})
}
