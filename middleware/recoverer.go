package middleware

import (
	"log"

	"github.com/lapitskyss/jsonrpc"
)

func Recovery() jsonrpc.MiddlewareFunc {
	return func(next jsonrpc.Handler) jsonrpc.Handler {
		return func(ctx *jsonrpc.RequestCtx) (_ jsonrpc.Result, err jsonrpc.Error) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.Println(rvr)

					err = jsonrpc.ErrInternalJSON()
				}
			}()

			return next(ctx)
		}
	}
}
