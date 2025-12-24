package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// Auth is a placeholder for future authentication/authorization middleware.
func Auth(r *ghttp.Request) {
	r.Middleware.Next()
}
