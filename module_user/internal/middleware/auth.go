package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// Auth 占位中间件，后续可接入 JWT/Session。
func Auth(r *ghttp.Request) {
	r.Middleware.Next()
}
