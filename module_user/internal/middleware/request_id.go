package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

// RequestID 注入或透传请求 ID。
func RequestID(r *ghttp.Request) {
	rid := r.Header.Get("X-Request-Id")
	if rid == "" {
		rid = guid.S()
	}
	r.SetCtxVar("request_id", rid)
	r.Response.Header().Set("X-Request-Id", rid)
	r.Middleware.Next()
}
