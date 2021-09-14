package jsonrpc

import (
	"bytes"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  Result           `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
}

type Result interface{}

// sendResponse send success JSON response
func sendResponse(ctx *fasthttp.RequestCtx, r []*Response) {
	for i := range r {
		if r[i].Error == nil && r[i].Result == nil {
			r[i].Result = json.RawMessage(nil)
		}
	}

	buf := &bytes.Buffer{}

	if len(r) == 1 {
		if err := json.NewEncoder(buf).Encode(r[0]); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
	} else if len(r) > 1 {
		if err := json.NewEncoder(buf).Encode(r); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBodyRaw(buf.Bytes())
}

// sendResponse send single error JSON response
func sendSingleErrorResponse(ctx *fasthttp.RequestCtx, error *Error) {
	buf := &bytes.Buffer{}
	response := &Response{
		Version: Version,
		ID:      nil,
		Error:   error,
	}

	if err := json.NewEncoder(buf).Encode(response); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBodyRaw(buf.Bytes())
}
