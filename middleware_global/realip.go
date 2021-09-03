package middleware_global

import (
	"net/http"

	"github.com/tomasen/realip"

	"github.com/lapitskyss/jsonrpc"
)

func RealIP() jsonrpc.MiddlewareGlobalFunc {
	return func(r *http.Request) *jsonrpc.Error {
		if rip := realip.RealIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		return nil
	}
}
